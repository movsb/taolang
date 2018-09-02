package main

type Parser struct {
	tokenizer *Tokenizer

	// TODO find a better way to do this
	skipSemicolon bool
}

func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{
		tokenizer: tokenizer,
	}
}

func (p *Parser) Parse() (program *Program, err error) {
	defer func() {
		err = toErr(recover())
	}()

	program = &Program{}
	for {
		stmt := p.parseGlobalStatement()
		if stmt == nil {
			break
		}
		program.stmts = append(program.stmts, stmt)
	}
	tk := p.next()
	if tk.typ != ttEOF {
		panicf("unexpected token: %v", tk)
	}

	return program, nil
}

func (p *Parser) expect(tt TokenType) Token {
	token := p.next()
	if token.typ != tt {
		panicf("unexpected token: %v", token)
	}
	return token
}

func (p *Parser) match(tts ...TokenType) (Token, bool) {
	tk := p.next()
	for _, tt := range tts {
		if tk.typ == tt {
			return tk, true
		}
	}
	p.undo(tk)
	return Token{}, false
}

func (p *Parser) next() Token {
	return p.tokenizer.Next()
}

func (p *Parser) undo(tk Token) {
	p.tokenizer.Undo(tk)
}

func (p *Parser) skip(tt TokenType) bool {
	if p.follow(tt) {
		p.next()
		return true
	}
	return false
}

func (p *Parser) peek() Token {
	return p.tokenizer.Peek()
}

func (p *Parser) follow(tt TokenType) bool {
	return p.peek().typ == tt
}

func (p *Parser) enter() {
	p.tokenizer.PushFrame()
}

func (p *Parser) leave(putback bool) {
	p.tokenizer.PopFrame(putback)
}

func (p *Parser) parseGlobalStatement() Statement {
	return p.parseStatement(true)
}

func (p *Parser) parseStatement(global bool) Statement {
	tk := p.peek()

	switch tk.typ {
	case ttLet:
		return p.parseVariableStatement()
	case ttFunction:
		// // don't know whether it is a function statement or function expression.
		// // but, if a function doesn't have a name, it must be function expression.
		// p.enter()
		fn := p.parseFunctionStatement()
		// if fn.expr.name == "" {
		// 	p.leave(true)
		// 	// TODO we should directly parse function expression statement since we knew it is.
		// 	stmt := p.parseExpressionStatement()
		// 	if stmt == nil {
		//
		// 		panic("anonymous function expression must be called immediately")
		// 	}
		// }
		// p.leave(false)
		return fn
	}

	if global {
		return nil
	}

	switch tk.typ {
	case ttReturn:
		return p.parseReturnStatement()
	case ttLeftBrace:
		// Notice: block statement skips parsing {} as object literal.
		return p.parseBlockStatement()
	case ttFor:
		return p.parseForStatement()
	case ttBreak:
		return p.parseBreakStatement()
	case ttIf:
		return p.parseIfStatement()
	}

	if stmt := p.parseExpressionStatement(); stmt != nil {
		return stmt
	}
	if stmt := p.parseAssignmentStatement(); stmt != nil {
		return stmt
	}

	return nil
}

func (p *Parser) parseVariableStatement() *VariableStatement {
	var v VariableStatement
	p.expect(ttLet)
	v.Name = p.expect(ttIdentifier).str
	if p.follow(ttAssign) {
		p.next()
		v.Expr = p.parseExpression()
	}
	p.expect(ttSemicolon)
	return &v
}

func (p *Parser) parseAssignmentStatement() (stmt *AssignmentStatement) {
	p.enter()
	defer func() {
		p.leave(stmt == nil)
	}()

	var as AssignmentStatement

	as.left = p.parseExpression()

	if _, ok := p.match(ttAssign); ok {
		as.right = p.parseExpression()
	} else if op, ok := p.match(
		ttStarStarAssign,
		ttStarAssign, ttDivideAssign, ttPercentAssign,
		ttPlusAssign, ttMinusAssign,
		ttLeftShiftAssign, ttRightShiftAssign,
		ttAndAssign, ttOrAssign, ttXorAssign,
	); ok {
		right := p.parseExpression()
		var binOp TokenType
		switch op.typ {
		case ttStarStarAssign:
			binOp = ttStarStar
		case ttStarAssign:
			binOp = ttMultiply
		case ttDivideAssign:
			binOp = ttDivision
		case ttPercentAssign:
			binOp = ttPercent
		case ttPlusAssign:
			binOp = ttAddition
		case ttMinusAssign:
			binOp = ttSubstraction
		case ttLeftShiftAssign:
			binOp = ttLeftShift
		case ttRightShiftAssign:
			binOp = ttRightShift
		case ttAndAssign:
			binOp = ttBitAnd
		case ttOrAssign:
			binOp = ttBitOr
		case ttXorAssign:
			binOp = ttBitXor
		default:
			panic("won't go here")
		}
		as.right = NewBinaryExpression(as.left, binOp, right)
	} else {
		return nil
	}

	if !p.skipSemicolon {
		// hah? !skip? skip?
		p.skip(ttSemicolon)
	} else {
		p.skipSemicolon = false
	}

	return &as
}

func (p *Parser) parseFunctionStatement() *FunctionStatement {
	var fn FunctionStatement
	p.expect(ttFunction)
	expr := p.parseFunctionExpression()
	fn.expr = expr
	return &fn
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	p.expect(ttReturn)
	expr := p.parseExpression()
	p.expect(ttSemicolon)
	return &ReturnStatement{
		expr: expr,
	}
}

func (p *Parser) parseExpressionStatement() (stmt *ExpressionStatement) {
	p.enter()
	defer func() {
		p.leave(stmt == nil)
	}()

	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	stmt = &ExpressionStatement{
		expr: expr,
	}
	if p.skip(ttSemicolon) {
		return stmt
	}
	return nil
}

func (p *Parser) parseBlockStatement() (stmt *BlockStatement) {
	if !p.follow(ttLeftBrace) {
		return nil
	}

	block := &BlockStatement{}
	p.expect(ttLeftBrace)
	for {
		stmt := p.parseStatement(false)
		if stmt == nil {
			break
		}
		block.stmts = append(block.stmts, stmt)
	}
	p.expect(ttRightBrace)

	return block
}

// All three parts of for-stmt (init, test, incr) can be omitted.
// If all these parts are omitted, the two semicolons can be omitted.
//   for [init]; [test]; [incr] {}
//   for {}
func (p *Parser) parseForStatement() *ForStatement {
	p.expect(ttFor)

	var fs ForStatement

	hasInit := false

	if p.follow(ttLet) {
		hasInit = true
		// TODO init can be assignment
		fs.init = p.parseVariableStatement()
	} else if p.follow(ttSemicolon) {
		hasInit = true
		p.expect(ttSemicolon)
	} else if !p.follow(ttLeftBrace) {
		fs.test = p.parseExpression()
		hasInit = false
	}

	if hasInit {
		// test
		if !p.follow(ttSemicolon) {
			fs.test = p.parseExpression()
			if fs.test == nil {
				panic("expr expected")
			} else {
				p.expect(ttSemicolon)
			}
		} else {
			p.next()
			// no test
		}
		// incr
		if !p.follow(ttLeftBrace) {
			// is expr?
			p.enter()
			fs.incr = p.parseExpression()
			if !p.follow(ttLeftBrace) {
				p.leave(true)
				fs.incr = nil
			} else {
				p.leave(false)
			}
			// is assignment?
			if fs.incr == nil {
				p.enter()
				p.skipSemicolon = true
				fs.incr = p.parseAssignmentStatement()
				if !p.follow(ttLeftBrace) {
					p.leave(true)
					fs.incr = nil
				} else {
					p.leave(false)
				}
			}
			if fs.incr == nil {
				panic("incr expected")
			}
		} else {
			// no incr
		}
	} else {
		if !p.follow(ttLeftBrace) {
			panic("for needs body")
		}
	}

	fs.block = p.parseBlockStatement()
	if fs.block == nil {
		panic("for needs body")
	}

	return &fs
}

func (p *Parser) parseBreakStatement() *BreakStatement {
	p.expect(ttBreak)
	p.expect(ttSemicolon)
	return &BreakStatement{}
}

func (p *Parser) parseIfStatement() *IfStatement {
	p.expect(ttIf)
	expr := p.parseExpression()
	ifBlock := p.parseBlockStatement()
	var elseBlock Statement
	switch p.peek().typ {
	case ttElse:
		p.expect(ttElse)
		switch p.peek().typ {
		case ttIf:
			elseBlock = p.parseIfStatement()
		case ttLeftBrace:
			elseBlock = p.parseBlockStatement()
		default:
			panic("else expect if or block to follow")
		}
	}
	return &IfStatement{
		cond:      expr,
		ifBlock:   ifBlock,
		elseBlock: elseBlock,
	}
}

func (p *Parser) parseExpression() Expression {
	return p.parseTernaryExpression()
}

func (p *Parser) parseTernaryExpression() Expression {
	var cond, left, right Expression
	cond = p.parseLogicalExpression()
	if _, ok := p.match(ttQuestion); ok {
		// Although we don't allow nested ternary expression, we parse it, we panic it, later.
		// left = p.parseLogicalExpression()
		left = p.parseExpression()
		p.expect(ttColon)
		right = p.parseExpression()
		errstr := " expression of ?: cannot be ?: (nested ?: is not allowed)"
		if _, ok := left.(*TernaryExpression); ok {
			panic("left" + errstr)
		}
		if _, ok := right.(*TernaryExpression); ok {
			panic("right" + errstr)
		}
		return NewTernaryExpression(cond, left, right)
	}
	return cond
}

func (p *Parser) parseLogicalExpression() Expression {
	left := p.parseBitwiseExpression()
	for {
		if op, ok := p.match(ttAndAnd, ttOrOr); ok {
			right := p.parseBitwiseExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
}

// Attention: &, |, ^ have the same precedence that doesn't match javascript.
func (p *Parser) parseBitwiseExpression() Expression {
	left := p.parseEqualityExpression()
	for {
		if op, ok := p.match(ttBitAnd, ttBitOr, ttBitXor); ok {
			right := p.parseEqualityExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
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
	left := p.parseShiftExpression()
	for {
		if op, ok := p.match(ttGreaterThan, ttGreaterThanOrEqual, ttLessThan, ttLessThanOrEqual); ok {
			right := p.parseShiftExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseShiftExpression() Expression {
	left := p.parseAdditionExpression()
	for {
		if op, ok := p.match(ttLeftShift, ttRightShift); ok {
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
	left := p.parseExponentiationExpression()
	for {
		if op, ok := p.match(ttMultiply, ttDivision, ttPercent); ok {
			right := p.parseExponentiationExpression()
			left = NewBinaryExpression(left, op.typ, right)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseExponentiationExpression() Expression {
	left := p.parseUnaryExpression()
	if op, ok := p.match(ttStarStar); ok {
		// ** is r2l association
		right := p.parseExponentiationExpression()
		return NewBinaryExpression(left, op.typ, right)
	}
	return left
}

func (p *Parser) parseUnaryExpression() Expression {
	if op, ok := p.match(ttNot, ttSubstraction, ttAddition); ok {
		right := p.parseUnaryExpression()
		return NewUnaryExpression(op.typ, right)
	}
	if op, ok := p.match(ttIncrement, ttDecrement); ok {
		right := p.parseUnaryExpression()
		return NewIncrementDecrementExpression(op, true, right)
	}
	return p.parsePostfixIncrementDecrementExpression()
}

func (p *Parser) parsePostfixIncrementDecrementExpression() Expression {
	left := p.parsePrimaryExpression()
	if left == nil {
		return left
	}
	if op, ok := p.match(ttIncrement, ttDecrement); ok {
		return NewIncrementDecrementExpression(op, false, left)
	}
	return left
}

func (p *Parser) parsePrimaryExpression() Expression {
	var expr Expression

	next := p.next()

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
		p.undo(next)
		if lambda := p.parseLambdaExpression(); lambda != nil {
			return lambda
		}
		p.next()
		expr = p.parseExpression()
		p.expect(ttRightParen)
	case ttIdentifier:
		if p.follow(ttLambda) {
			p.undo(next)
			lambda := p.parseLambdaExpression()
			if lambda == nil {
				panic("bad lambda expression")
			}
			return lambda
		}
		expr = ValueFromVariable(next.str)
	case ttFunction:
		expr = p.parseFunctionExpression()
	case ttLeftBrace:
		p.undo(next)
		expr = p.parseObjectExpression()
	case ttLeftBracket:
		p.undo(next)
		expr = p.parseArrayExpression()
	default:
		p.undo(next)
		return nil
	}

	for {
		if index := p.parseIndexExpression(expr); index != nil {
			expr = index
			continue
		}
		if call := p.parseCallExpression(expr); call != nil {
			expr = call
			continue
		}
		break
	}

	return expr
}

func (p *Parser) parseLambdaExpression() (expr *FunctionExpression) {
	params := &Parameters{}

	if _, ok := p.match(ttLeftParen); ok {
		for {
			params.PutParam(p.expect(ttIdentifier).str)
			if p.follow(ttRightParen) {
				break
			}
			if !p.skip(ttComma) {
				return nil
			}
		}
		p.expect(ttRightParen)
	} else {
		params.PutParam(p.expect(ttIdentifier).str)
	}

	p.expect(ttLambda)

	var block *BlockStatement

	if p.follow(ttLeftBrace) {
		block = p.parseBlockStatement()
	} else {
		expr := p.parseExpression()
		if expr == nil {
			panic("cannot parse lambda body expr")
		}
		ret := NewReturnStatement(expr)
		block = NewBlockStatement(ret)
	}

	return &FunctionExpression{
		params: params,
		block:  block,
	}
}

func (p *Parser) parseIndexExpression(left Expression) (expr Expression) {
	switch token := p.next(); token.typ {
	case ttDot:
		if ident := p.next(); ident.typ == ttIdentifier {
			return &IndexExpression{
				indexable: left,
				key:       ValueFromString(ident.str),
			}
		}
		// panic("unknown token after expr")
	case ttLeftBracket:
		keyExpr := p.parseExpression()
		if keyExpr == nil {
			return nil
		}
		if bracket := p.next(); bracket.typ != ttRightBracket {
			return nil
		}
		return &IndexExpression{
			indexable: left,
			key:       keyExpr,
		}
	default:
		p.undo(token)
	}
	return nil
}

func (p *Parser) parseCallExpression(left Expression) Expression {
	if paren := p.next(); paren.typ != ttLeftParen {
		p.undo(paren)
		return nil
	}

	call := CallExpression{
		Callable: left,
		Args:     &Arguments{},
	}

	for {
		arg := p.parseExpression()
		if arg == nil {
			break
		}
		call.Args.PutArgument(arg)
		if !p.skip(ttComma) {
			break
		}
	}

	p.expect(ttRightParen)

	return &call
}

func (p *Parser) parseFunctionExpression() *FunctionExpression {
	var name string
	var block *BlockStatement
	params := &Parameters{}

	if p.follow(ttIdentifier) {
		name = p.next().str
	}

	p.expect(ttLeftParen)
	for {
		tk := p.next()
		if tk.typ == ttIdentifier {
			params.PutParam(tk.str)
		} else if tk.typ == ttComma {
			continue
		} else if tk.typ == ttRightParen {
			p.undo(tk)
			break
		}
	}
	p.expect(ttRightParen)

	if !p.follow(ttLeftBrace) {
		panic("function needs a body")
	}

	block = p.parseBlockStatement()

	return &FunctionExpression{
		name:   name,
		params: params,
		block:  block,
	}
}

func (p *Parser) parseObjectExpression() Expression {
	objexpr := NewObjectExpression()

	var key string
	var expr Expression

	p.expect(ttLeftBrace)

	for {
		token := p.next()
		if token.typ == ttRightBrace {
			p.undo(token)
			break
		}

		switch token.typ {
		case ttString:
			key = token.str
		case ttIdentifier:
			key = token.str
		default:
			panic("unsupported key type")
		}

		p.expect(ttColon)

		expr = p.parseExpression()
		objexpr.props[key] = expr

		p.skip(ttComma)

		if p.follow(ttRightBrace) {
			break
		}
	}

	p.expect(ttRightBrace)

	return objexpr
}

func (p *Parser) parseArrayExpression() Expression {
	arrExpr := NewArrayExpression()

	p.expect(ttLeftBracket)

	for {
		elem := p.parseExpression()
		if elem == nil {
			break
		}
		arrExpr.elements = append(arrExpr.elements, elem)
		p.skip(ttComma)
		if p.follow(ttRightBracket) {
			break
		}
	}

	p.expect(ttRightBracket)

	return arrExpr
}
