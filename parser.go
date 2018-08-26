package main

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{
		tokenizer: tokenizer,
	}
}

func (p *Parser) Parse() *Program {
	program := Program{}
	for {
		stmt := p.parseGlobalStatement()
		if stmt == nil {
			break
		}
		program.stmts = append(program.stmts, stmt)
	}
	tk := p.tokenizer.Next()
	if tk.typ != ttEOF {
		panic("unexpected statement")
	}
	return &program
}

func (p *Parser) expect(tt TokenType) Token {
	token := p.tokenizer.Next()
	if token.typ != tt {
		panic("unexpected token type")
	}
	return token
}

func (p *Parser) match(tts ...TokenType) (Token, bool) {
	tk := p.tokenizer.Next()
	for _, tt := range tts {
		if tk.typ == tt {
			return tk, true
		}
	}
	p.tokenizer.Undo(tk)
	return Token{}, false
}

func (p *Parser) parseGlobalStatement() Statement {
	return p.parseStatement(true)
}

func (p *Parser) parseStatement(global bool) Statement {
	tk := p.tokenizer.Next()

	switch tk.typ {
	case ttLet:
		return p.parseVariableDefinitionStatement()
	case ttFunction:
		stmt := p.parseFunctionDefinitionStatement()
		fn := stmt.(*FunctionDefinitionStatement)
		if fn.name == "" {
			panic("function statement must have function name")
		}
		return stmt
	}

	if global {
		p.tokenizer.Undo(tk)
		return nil
	}

	switch tk.typ {
	case ttReturn:
		return p.parseReturnStatement()
	}

	p.tokenizer.Undo(tk)

	if stmt := p.parseExpressionStatement(); stmt != nil {
		return stmt
	}
	if stmt := p.parseAssignmentStatement(); stmt != nil {
		return stmt
	}

	return nil
}

func (p *Parser) parseExpression() Expression {
	return p.parseEqualityExpression()
}

func (p *Parser) parseVariableDefinitionStatement() Statement {
	var v VariableDefinitionStatement
	v.Name = p.expect(ttIdentifier).str
	if p.tokenizer.Peek().typ == ttAssign {
		p.tokenizer.Next()
		v.Expr = p.parseExpression()
	}
	p.expect(ttSemicolon)
	return &v
}

func (p *Parser) parseAssignmentStatement() (stmt Statement) {
	p.tokenizer.PushFrame()
	defer func() {
		p.tokenizer.PopFrame(stmt == nil)
	}()

	var as VariableAssignmentStatement
	name := p.tokenizer.Next()
	if name.typ != ttIdentifier &&
		name.typ != ttBoolean && // these two are predeclared constants
		name.typ != ttNil {
		p.tokenizer.Undo(name)
		return nil
	}
	as.Name = name.str

	assign := p.tokenizer.Next()
	if assign.typ != ttAssign {
		p.tokenizer.Undo(assign)
		p.tokenizer.Undo(name)
		return nil
	}

	as.Expr = p.parseExpression()

	semi := p.tokenizer.Next()
	if semi.typ != ttSemicolon {
		p.tokenizer.Undo(semi)
		return nil
	}

	if name.typ == ttBoolean || name.typ == ttNil {
		panic("predeclared constants cannot be assigned")
	}

	return &as
}

func (p *Parser) parseFunctionDefinitionStatement() Statement {
	var fn FunctionDefinitionStatement
	expr := p.parseFunctionDefinitionExpression().(*FunctionDefinitionExpression)
	fn.name = expr.name
	fn.expr = expr
	return &fn
}

func (p *Parser) parseReturnStatement() Statement {
	expr := p.parseExpression()
	p.expect(ttSemicolon)
	return &ReturnStatement{
		expr: expr,
	}
}

func (p *Parser) parseExpressionStatement() (stmt Statement) {
	p.tokenizer.PushFrame()
	defer func() {
		p.tokenizer.PopFrame(stmt == nil)
	}()

	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	stmt = &ExpressionStatement{
		expr: expr,
	}
	semi := p.tokenizer.Next()
	if semi.typ == ttSemicolon {
		return stmt
	}
	p.tokenizer.Undo(semi)
	return nil
}

func (p *Parser) parseEqualityExpression() Expression {
	left := p.parseComparisonExpression()
	for {
		if op, ok := p.match(ttEqual, ttNotEqual); ok {
			right := p.parseComparisonExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseComparisonExpression() Expression {
	left := p.parseAdditionExpression()
	for {
		if op, ok := p.match(ttGreaterThan, ttGreaterThanOrEqual, ttLessThan, ttLessThanOrEqual); ok {
			right := p.parseAdditionExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseAdditionExpression() Expression {
	left := p.parseMultiplicationExpression()
	for {
		if op, ok := p.match(ttAddition, ttSubstraction); ok {
			right := p.parseMultiplicationExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseMultiplicationExpression() Expression {
	left := p.parseUnaryExpression()
	for {
		if op, ok := p.match(ttMultiply, ttDivision); ok {
			right := p.parseUnaryExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseUnaryExpression() Expression {
	if op, ok := p.match(ttNot, ttSubstraction); ok {
		right := p.parseUnaryExpression()
		return NewUnaryExpression(op.typ, right)
	}
	return p.parsePrimaryExpression()
}

func (p *Parser) parsePrimaryExpression() Expression {
	var expr Expression
	next := p.tokenizer.Next()

	switch next.typ {
	case ttNil:
		expr = ValueFromNil()
	case ttBoolean:
		expr = ValueFromBoolean(next.str == "true")
	case ttNumber:
		expr = ValueFromNumber(next.num)
	case ttString:
		expr = ValueFromString(next.str)
	case ttLeftParen:
		expr = p.parseExpression()
		p.expect(ttRightParen)
	case ttIdentifier:
		expr = ValueFromVariable(next.str)
	case ttFunction:
		expr = p.parseFunctionDefinitionExpression()
	default:
		p.tokenizer.Undo(next)
		return nil
	}

	if call := p.parseCallExpression(); call != nil {
		callExpr := call.(*CallExpression)
		callExpr.Callable = expr
		return call
	}

	return expr
}

func (p *Parser) parseCallExpression() Expression {
	if paren := p.tokenizer.Next(); paren.typ != ttLeftParen {
		p.tokenizer.Undo(paren)
		return nil
	}

	call := CallExpression{}
	call.Args = &Arguments{}

	for {
		arg := p.parseExpression()
		if arg == nil {
			break
		}
		call.Args.PutArgument(arg)
		if comma := p.tokenizer.Next(); comma.typ != ttComma {
			p.tokenizer.Undo(comma)
			break
		}
	}

	p.expect(ttRightParen)

	return &call
}

func (p *Parser) parseFunctionDefinitionExpression() Expression {
	var name string
	params := &Parameters{}
	stmts := []Statement{}

	if p.tokenizer.Peek().typ == ttIdentifier {
		name = p.tokenizer.Next().str
	}

	p.expect(ttLeftParen)
	for {
		tk := p.tokenizer.Next()
		if tk.typ == ttIdentifier {
			params.PutParam(tk.str)
		} else if tk.typ == ttComma {
			continue
		} else if tk.typ == ttRightParen {
			p.tokenizer.Undo(tk)
			break
		}
	}
	p.expect(ttRightParen)

	p.expect(ttLeftBrace)
	for {
		stmt := p.parseStatement(false)
		if stmt == nil {
			break
		}
		stmts = append(stmts, stmt)
	}
	p.expect(ttRightBrace)

	return &FunctionDefinitionExpression{
		name:   name,
		params: params,
		stmts:  stmts,
	}
}
