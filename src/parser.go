package main

// Precedence is operator precedence table.
type Precedence uint

const (
	_ Precedence = iota

	precAssignment
	precConditional
	precLogicalOr
	precLogicalAnd

	precBitwiseAnd,
	precBitwiseOr,
	precBitwiseXor,
	precBitwiseAndNot = iota, iota, iota, iota

	precEquality = iota
	precComparison
	precBitwiseShift
	precAddition
	precMultiplication
	precExponentiation

	precLogicalNot,
	precBitwiseNot,
	precUnaryPlus,
	precUnaryNegation = iota, iota, iota, iota

	precPrefixIncrement,
	precPrefixDecrement = iota, iota

	precPostfixIncrement,
	precPostfixDecrement = iota, iota

	precIndexing,
	precNew,
	precCall = iota, iota, iota

	precHighest = iota
)

// Precedence table.
// Assignment is not in here.
//
// Operators that can be either unary or binary:
//   +, -, ^
// These operators map to two precedence.
//
// Increment / Decrement can be either:
// 		prefix: ++a
// or:
//		postfix: a--
//
var precedenceTable = map[TokenType]Precedence{
	ttQuestion:           precConditional,
	ttLogicalNot:         precLogicalNot,
	ttLogicalOr:          precLogicalOr,
	ttLogicalAnd:         precLogicalAnd,
	ttBitAnd:             precBitwiseAnd,
	ttBitOr:              precBitwiseOr,
	ttBitXor:             precBitwiseXor,
	ttBitAndNot:          precBitwiseAndNot,
	ttEqual:              precEquality,
	ttNotEqual:           precEquality,
	ttGreaterThan:        precComparison,
	ttGreaterThanOrEqual: precComparison,
	ttLessThan:           precComparison,
	ttLessThanOrEqual:    precComparison,
	ttLeftShift:          precBitwiseShift,
	ttRightShift:         precBitwiseShift,
	ttAddition:           precAddition,
	ttSubtraction:        precAddition,
	ttMultiply:           precMultiplication,
	ttDivision:           precMultiplication,
	ttPercent:            precMultiplication,
	ttStarStar:           precExponentiation,
	ttIncrement:          precPrefixIncrement,
	ttDecrement:          precPrefixDecrement,
	ttLeftBracket:        precIndexing,
	ttDot:                precIndexing,
	ttNew:                precNew,
	ttLeftParen:          precCall,
}

// Parser parses tokens into AbstractSyntaxTree.
type Parser struct {
	tokenizer *Tokenizer

	// how many breakable states the parser is in.
	// a positive breakCount means we can break breakCount times.
	breakCount uint
}

// NewParser news a parser.
func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{
		tokenizer:  tokenizer,
		breakCount: 0,
	}
}

// Parse does parse the input tokens.
func (p *Parser) Parse() (program *Program, err interface{}) {
	defer func() {
		err = recover()
	}()

	program = &Program{}
	for {
		if p.follow(ttEOF) {
			break
		}
		stmt := p.parseStatement(true)
		program.stmts = append(program.stmts, stmt)
	}
	tk := p.next()
	if tk.typ != ttEOF {
		panic(NewSyntaxError("unexpected token: %v", tk))
	}

	return program, nil
}

func (p *Parser) expect(tt TokenType) Token {
	token := p.next()
	if token.typ != tt {
		exp := Token{typ: tt}.String()
		switch tt {
		case ttIdentifier:
			exp = "`identifier'"
		}
		panic(NewSyntaxError("unexpected token: %v (expect: %v)", token, exp))
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

func (p *Parser) getPrecedence(op TokenType) (prec Precedence, ok bool) {
	if op >= ttAssign && op <= ttAndNotAssign {
		return precAssignment, true
	}
	if prec, ok = precedenceTable[op]; ok {
		return
	}
	return
}

func (p *Parser) parseStatement(global bool) Statement {
	tk := p.peek()

	switch tk.typ {
	case ttLet:
		return p.parseLetStatement()
	case ttFunction:
		return p.parseFunctionStatement()
	case ttSemicolon:
		p.next()
		return &EmptyStatement{}
	}

	if global {
		panic(NewSyntaxError("non-global statement"))
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
		if p.breakCount <= 0 {
			panic(NewSyntaxError("break statement must be in for-loop or switch"))
		}
		return p.parseBreakStatement()
	case ttIf:
		return p.parseIfStatement()
	case ttSwitch:
		return p.parseSwitchStatement()
	case ttTao:
		return p.parseTaoStatement()
	}

	// At last, try to parse all another statements we've known.
	//   - expr;                    expression statement
	{
		expr := p.parseExpression(precAssignment)
		// it is an expression statement
		if _, ok := p.match(ttSemicolon); ok {
			return &ExpressionStatement{
				expr: expr,
			}
		}
		p.expect(ttSemicolon)
	}

	panic(NewSyntaxError("unknown statement at line: %d", tk.line))
}

func (p *Parser) parseLetStatement() *LetStatement {
	var l LetStatement
	p.expect(ttLet)
	l.Name = p.expect(ttIdentifier).str
	if p.follow(ttAssign) {
		p.next()
		l.Expr = p.parseExpression(precConditional)
	}
	p.expect(ttSemicolon)
	return &l
}

func (p *Parser) parseFunctionStatement() *FunctionStatement {
	var fn FunctionStatement
	fn.expr = p.parseFunctionExpression()
	return &fn
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	p.expect(ttReturn)
	ret := &ReturnStatement{}
	if !p.follow(ttSemicolon) {
		ret.expr = p.parseExpression(precConditional)
	}
	p.expect(ttSemicolon)
	return ret
}

func (p *Parser) parseBlockStatement() (stmt *BlockStatement) {
	block := &BlockStatement{}
	p.expect(ttLeftBrace)
	for {
		if p.follow(ttRightBrace) {
			break
		}
		stmt := p.parseStatement(false)
		block.stmts = append(block.stmts, stmt)
	}
	p.expect(ttRightBrace)

	return block
}

// All three parts of for-stmt (init, test, incr) can be omitted.
// If all these parts are omitted, the two semicolons can be omitted.
//   for [init]; [test]; [incr] {}
//   for expr {}
//   for {}
func (p *Parser) parseForStatement() *ForStatement {
	p.expect(ttFor)

	var fs ForStatement

	hasInit := false

	if p.follow(ttLet) {
		hasInit = true
		// TODO init can be assignment
		fs.init = p.parseLetStatement()
	} else if p.follow(ttSemicolon) {
		hasInit = true
		p.expect(ttSemicolon)
	} else if !p.follow(ttLeftBrace) {
		fs.test = p.parseExpression(precConditional)
		hasInit = false
	}

	if hasInit {
		// test
		if !p.follow(ttSemicolon) {
			fs.test = p.parseExpression(precConditional)
			p.expect(ttSemicolon)
		} else {
			p.next()
			// no test
		}
		// incr
		if !p.follow(ttLeftBrace) {
			fs.incr = p.parseExpression(precAssignment)
		} else {
			// no incr
		}
	} else {
		if !p.follow(ttLeftBrace) {
			panic(NewSyntaxError("for needs body"))
		}
	}

	p.breakCount++

	fs.block = p.parseBlockStatement()
	if fs.block == nil {
		panic(NewSyntaxError("for needs body"))
	}

	p.breakCount--

	return &fs
}

func (p *Parser) parseBreakStatement() *BreakStatement {
	p.expect(ttBreak)
	p.expect(ttSemicolon)
	return &BreakStatement{}
}

func (p *Parser) parseIfStatement() *IfStatement {
	p.expect(ttIf)
	expr := p.parseExpression(precConditional)
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
			panic(NewSyntaxError("else expect if or block to follow"))
		}
	}
	return &IfStatement{
		cond:      expr,
		ifBlock:   ifBlock,
		elseBlock: elseBlock,
	}
}

func (p *Parser) parseSwitchStatement() *SwitchStatement {
	p.expect(ttSwitch)
	ss := SwitchStatement{}
	ss.cond = p.parseExpression(precConditional)
	p.expect(ttLeftBrace)
	// empty cases
	if p.follow(ttRightBrace) {
		p.next()
		return &ss
	}
	// parse cases
	for {
		var group *CaseGroup
		switch esac := p.next(); esac.typ {
		case ttDefault:
			if ss.def != nil {
				panic(NewSyntaxError("duplicate default"))
			}
			p.expect(ttColon)
			group = &CaseGroup{}
			ss.def = group
		case ttCase:
			group = &CaseGroup{}
			ss.cases = append(ss.cases, group)
			for {
				expr := p.parseExpression(precConditional)
				group.cases = append(group.cases, expr)
				p.skip(ttComma)
				if p.follow(ttColon) {
					p.next()
					break
				}
			}
		default:
			panic(NewSyntaxError("unexpected token: %v", esac))
		}
		// no need to be real block.
		// but we construct a block since statements are executed in scope.
		group.block = &BlockStatement{}
		// we can break in switch statement
		p.breakCount++
		for {
			if tk, ok := p.match(ttCase, ttDefault, ttRightBrace); ok {
				p.undo(tk)
				break
			}
			stmt := p.parseStatement(false)
			group.block.stmts = append(group.block.stmts, stmt)
		}
		p.breakCount--
		if p.follow(ttRightBrace) {
			p.next()
			break
		}
	}
	return &ss
}

func (p *Parser) parseTaoStatement() *TaoStatement {
	p.expect(ttTao)
	ts := TaoStatement{}
	expr := p.parseExpression(precConditional)
	if _, ok := expr.(*CallExpression); !ok {
		panic(NewSyntaxError("expression in tao must be function call"))
	}
	p.expect(ttSemicolon)
	ts.call = expr.(*CallExpression)
	return &ts
}

func (p *Parser) parseExpression(prec Precedence) Expression {
	var left Expression
	peek := p.peek()
	if prec <= precUnaryPlus {
		switch peek.typ {
		case ttLogicalNot, ttBitXor, ttAddition, ttSubtraction:
			p.next()
			right := p.parseExpression(precUnaryPlus)
			left = NewUnaryExpression(peek.typ, right)
		}
	}
	if prec <= precPrefixIncrement {
		switch peek.typ {
		case ttIncrement, ttDecrement:
			p.next()
			right := p.parseExpression(precPrefixIncrement)
			left = NewIncrementDecrementExpression(peek.typ, true, right)
		}
	}
	if prec <= precNew {
		if peek.typ == ttNew {
			left = p.parseNewExpression()
		}
	}
	if left == nil {
		left = p.parsePrimaryExpression()
	}

	for {
		op := p.next()

		nextPrec, ok := p.getPrecedence(op.typ)
		if !ok || nextPrec < prec {
			p.undo(op)
			break
		}

		if op.typ >= ttAssign && op.typ < ttQuestion {
			// return: cannot use assignment as expression
			return p.parseAssignmentExpression(left, op.typ)
		}

		switch op.typ {
		case ttQuestion:
			left = p.parseTernaryExpression(left)
			continue
		case ttIncrement, ttDecrement:
			left = NewIncrementDecrementExpression(op.typ, false, left)
			continue
		case ttLeftParen:
			p.undo(op)
			left = p.parseCallExpression(left)
			continue
		case ttLeftBracket, ttDot:
			p.undo(op)
			left = p.parseIndexExpression(left)
			continue
		}

		var right Expression

		switch op.typ {
		case ttLogicalOr:
			right = p.parseExpression(precLogicalAnd)
		case ttLogicalAnd:
			right = p.parseExpression(precBitwiseAnd)
		case ttBitAnd, ttBitOr, ttBitXor, ttBitAndNot:
			right = p.parseExpression(precEquality)
		case ttEqual, ttNotEqual:
			right = p.parseExpression(precComparison)
		case ttGreaterThan, ttGreaterThanOrEqual, ttLessThan, ttLessThanOrEqual:
			right = p.parseExpression(precBitwiseShift)
		case ttLeftShift, ttRightShift:
			right = p.parseExpression(precAddition)
		case ttAddition, ttSubtraction:
			right = p.parseExpression(precMultiplication)
		case ttMultiply, ttDivision, ttPercent:
			right = p.parseExpression(precExponentiation)
		case ttStarStar:
			right = p.parseExpression(precIndexing)
		}

		if right != nil {
			left = NewBinaryExpression(left, op.typ, right)
			continue
		}

		panic(NewSyntaxError("unknown operator: %v", op))
	}

	return left
}

func (p *Parser) parseTernaryExpression(cond Expression) Expression {
	var left, right Expression
	// Although we don't allow nested ternary expression, we parse it, we panic it, later.
	// left = p.parseLogicalExpression()
	left = p.parseExpression(precConditional)
	p.expect(ttColon)
	right = p.parseExpression(precConditional)
	if _, ok := left.(*TernaryExpression); ok {
		panic(NewSyntaxError("nested `?:' is not allowed"))
	}
	if _, ok := right.(*TernaryExpression); ok {
		panic(NewSyntaxError("nested `?:' is not allowed"))
	}
	return NewTernaryExpression(cond, left, right)
}

func (p *Parser) parseAssignmentExpression(left Expression, op TokenType) Expression {
	var ae AssignmentExpression
	ae.left = left

	// ttQuestion: disable continuous assignment style
	ae.right = p.parseExpression(precConditional)

	if op == ttAssign {
		return &ae
	}

	var binOp TokenType
	switch op {
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
		binOp = ttSubtraction
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
	case ttAndNotAssign:
		binOp = ttBitAndNot
	default:
		panic("won't go here")
	}

	ae.right = NewBinaryExpression(ae.left, binOp, ae.right)
	return &ae
}

func (p *Parser) parsePrimaryExpression() Expression {
	var expr Expression

	switch next := p.next(); next.typ {
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
		if lambda := p.tryParseLambdaExpression(false); lambda != nil {
			return lambda
		}
		p.next()
		expr = p.parseExpression(precConditional)
		p.expect(ttRightParen)
	case ttIdentifier:
		if p.follow(ttLambda) {
			p.undo(next)
			return p.tryParseLambdaExpression(true)
		}
		expr = ValueFromVariable(next.str)
	case ttFunction:
		p.undo(next)
		expr = p.parseFunctionExpression()
	case ttLeftBrace:
		p.undo(next)
		expr = p.parseObjectExpression()
	case ttLeftBracket:
		p.undo(next)
		expr = p.parseArrayExpression()
	default:
		panic(NewSyntaxError("unexpected token: %v", next))
	}
	return expr
}

func (p *Parser) parseNewExpression() *NewExpression {
	ne := &NewExpression{}
	p.expect(ttNew)
	ne.Type = p.expect(ttIdentifier).str
	ne.Args = &Arguments{}
	p.expect(ttLeftParen)
	if !p.follow(ttRightParen) {
		for {
			arg := p.parseExpression(precConditional)
			ne.Args.PutArgument(arg)
			sep := p.next()
			if sep.typ == ttComma {
				continue
			} else if sep.typ == ttRightParen {
				p.undo(sep)
				break
			} else {
				panic(NewSyntaxError("unexpected token: %v", sep))
			}
		}
	}
	p.expect(ttRightParen)
	return ne
}

// tryParseLambdaExpression tries to parse a lambda expression.
// If must, and no lambda can be parsed, it will panic.
// otherwise, it simply returns nil and leaves the tokens unchanged.
func (p *Parser) tryParseLambdaExpression(must bool) (expr *FunctionExpression) {
	defer func() {
		if must && expr == nil {
			panic(NewSyntaxError("bad lambda expression"))
		}
	}()

	p.enter()

	params := &Parameters{}

	if _, ok := p.match(ttLeftParen); ok {
		for {
			if p.follow(ttIdentifier) {
				params.PutParam(p.expect(ttIdentifier).str)
			} else {
				if p.follow(ttRightParen) {
					break
				} else {
					p.leave(true)
					return nil
				}
			}
			if p.skip(ttComma) {
				continue
			}
		}
		if _, ok := p.match(ttRightParen); !ok {
			p.leave(true)
			return nil
		}
	} else {
		if !p.follow(ttIdentifier) {
			p.leave(true)
			return nil
		}
		params.PutParam(p.next().str)
	}

	if !p.follow(ttLambda) {
		p.leave(true)
		return nil
	}

	p.leave(false)
	p.next()

	var body *BlockStatement

	if p.follow(ttLeftBrace) {
		body = p.parseBlockStatement()
	} else {
		expr := p.parseExpression(precConditional)
		ret := NewReturnStatement(expr)
		body = NewBlockStatement(ret)
	}

	// a lambda expression is just an anonymous function
	return &FunctionExpression{
		name:   "",
		params: params,
		body:   body,
	}
}

func (p *Parser) parseIndexExpression(left Expression) (expr Expression) {
	switch token := p.next(); token.typ {
	case ttDot:
		key := p.next()
		if key.typ == ttIdentifier {
			return &IndexExpression{
				indexable: left,
				key:       ValueFromString(key.str),
			}
		}
		panic(NewSyntaxError("unexpected token: %v", key))
	case ttLeftBracket:
		keyExpr := p.parseExpression(precConditional)
		if bracket := p.next(); bracket.typ != ttRightBracket {
			return nil
		}
		return &IndexExpression{
			indexable: left,
			key:       keyExpr,
		}
	default:
		panic("won't go here")
	}
}

func (p *Parser) parseCallExpression(left Expression) Expression {
	call := CallExpression{
		Callable: left,
	}

	p.expect(ttLeftParen)

	if !p.follow(ttRightParen) {
		for {
			arg := p.parseExpression(precConditional)
			call.Args.PutArgument(arg)
			sep := p.next()
			if sep.typ == ttComma {
				continue
			} else if sep.typ == ttRightParen {
				p.undo(sep)
				break
			} else {
				panic(NewSyntaxError("unexpected token: %v", sep))
			}
		}
	}

	p.expect(ttRightParen)

	return &call
}

func (p *Parser) parseFunctionExpression() *FunctionExpression {
	var fn FunctionExpression
	fn.params = &Parameters{}

	p.expect(ttFunction)

	if p.follow(ttIdentifier) {
		fn.name = p.next().str
	}

	p.expect(ttLeftParen)
	if !p.follow(ttRightParen) {
		for {
			name := p.expect(ttIdentifier).str
			fn.params.PutParam(name)
			sep := p.next()
			if sep.typ == ttComma {
				continue
			} else if sep.typ == ttRightParen {
				p.undo(sep)
				break
			} else {
				panic(NewSyntaxError("unexpected token: %v", sep))
			}
		}
	}
	p.expect(ttRightParen)

	if !p.follow(ttLeftBrace) {
		panic(NewSyntaxError("function needs a body"))
	}

	saveBreakCount := p.breakCount
	p.breakCount = 0

	fn.body = p.parseBlockStatement()

	p.breakCount = saveBreakCount

	return &fn
}

func (p *Parser) parseObjectExpression() Expression {
	objexpr := NewObjectExpression()

	p.expect(ttLeftBrace)

	for {
		if p.follow(ttRightBrace) {
			break
		}

		var key string
		var expr Expression

		switch token := p.next(); token.typ {
		case ttString:
			key = token.str
		case ttIdentifier:
			key = token.str
		default:
			panic(NewTypeError("unsupported key type"))
		}

		p.expect(ttColon)

		expr = p.parseExpression(precConditional)
		objexpr.props[key] = expr

		// allow last comma
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
		if p.follow(ttRightBracket) {
			break
		}

		elem := p.parseExpression(precConditional)
		arrExpr.elements = append(arrExpr.elements, elem)
		// allow last comma
		p.skip(ttComma)
		if p.follow(ttRightBracket) {
			break
		}
	}

	p.expect(ttRightBracket)

	return arrExpr
}
