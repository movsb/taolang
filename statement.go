package main

type Statement interface {
	Execute(ctx *Context)
}

type Returner interface {
	Return() (Value, bool)
}

type Breaker interface {
	Break() bool
}

type VariableStatement struct {
	Name string
	Expr Expression
}

func (v *VariableStatement) Execute(ctx *Context) {
	value := ValueFromNil()
	if v.Expr != nil {
		value = v.Expr.Evaluate(ctx)
	}
	ctx.AddValue(v.Name, value)
}

type VariableAssignmentStatement struct {
	Name string
	Expr Expression
}

func (v *VariableAssignmentStatement) Execute(ctx *Context) {
	ctx.SetValue(v.Name, v.Expr.Evaluate(ctx))
}

type FunctionStatement struct {
	name string
	expr *FunctionExpression
}

func (f *FunctionStatement) Execute(ctx *Context) {
	ctx.AddValue(f.name, ValueFromFunction(f.expr))
}

type ReturnStatement struct {
	expr  Expression
	value Value
}

func NewReturnStatement(expr Expression) *ReturnStatement {
	return &ReturnStatement{
		expr: expr,
	}
}

func (r *ReturnStatement) Execute(ctx *Context) {
	r.value = r.expr.Evaluate(ctx)
}

func (r *ReturnStatement) Return() (Value, bool) {
	return r.value, true
}

type BlockStatement struct {
	retValue Value
	broke    bool
	stmts    []Statement
}

func NewBlockStatement(stmts ...Statement) *BlockStatement {
	b := &BlockStatement{}
	for _, stmt := range stmts {
		b.stmts = append(b.stmts, stmt)
	}
	return b
}

func (b *BlockStatement) Return() (Value, bool) {
	return b.retValue, b.retValue.defined()
}

func (b *BlockStatement) Break() bool {
	return b.broke
}

func (b *BlockStatement) Execute(ctx *Context) {
	for _, stmt := range b.stmts {
		switch typed := stmt.(type) {
		case *BlockStatement:
			newCtx := NewContext(ctx)
			typed.Execute(newCtx)
		default:
			typed.Execute(ctx)
		}
		if returner, ok := stmt.(Returner); ok {
			if ret, ok := returner.Return(); ok {
				b.retValue = ret
				return
			}
		}
		if breaker, ok := stmt.(Breaker); ok {
			if breaker.Break() {
				b.broke = true
				break
			}
		}
	}
}

type ExpressionStatement struct {
	expr Expression
}

func NewExpressionStatement(expr Expression) *ExpressionStatement {
	return &ExpressionStatement{
		expr: expr,
	}
}

func (r *ExpressionStatement) Execute(ctx *Context) {
	value := r.expr.Evaluate(ctx)
	_ = value // drop expr value
}

type WhileStatement struct {
	expr     Expression
	block    *BlockStatement
	retValue Value
}

func (w *WhileStatement) Return() (Value, bool) {
	return w.retValue, w.retValue.defined()
}

func (w *WhileStatement) Execute(ctx *Context) {
	for {
		cond := w.expr.Evaluate(ctx)
		if !cond.Truth(ctx) {
			break
		}
		newCtx := NewContext(ctx)
		w.block.Execute(newCtx)
		if value, ok := w.block.Return(); ok {
			w.retValue = value
			break
		}
		if w.block.Break() {
			break
		}
	}
}

type BreakStatement struct {
}

func (b *BreakStatement) Break() bool {
	return true
}

func (b *BreakStatement) Execute(ctx *Context) {

}

type IfStatement struct {
	cond      Expression
	ifBlock   *BlockStatement
	elseBlock Statement // if or block
	retValue  Value
	broke     bool
}

func (i *IfStatement) Return() (Value, bool) {
	return i.retValue, i.retValue.defined()
}

func (i *IfStatement) Break() bool {
	return i.broke
}

func (i *IfStatement) Execute(ctx *Context) {
	cond := i.cond.Evaluate(ctx)
	if cond.Truth(ctx) {
		newCtx := NewContext(ctx)
		i.ifBlock.Execute(newCtx)
		if ret, ok := i.ifBlock.Return(); ok {
			i.retValue = ret
			return
		}
		if broke := i.ifBlock.Break(); broke {
			i.broke = true
			return
		}
	} else {
		var stmt Statement
		switch typed := i.elseBlock.(type) {
		case nil:
			return
		case *IfStatement:
			stmt = typed
			typed.Execute(ctx)
		case *BlockStatement:
			stmt = typed
			newCtx := NewContext(ctx)
			typed.Execute(newCtx)
		default:
			panic("bad else stmt")
		}
		if ret, ok := stmt.(Returner); ok {
			if value, ok := ret.Return(); ok {
				i.retValue = value
				return
			}
		}
		if brk, ok := stmt.(Breaker); ok {
			if brk.Break() {
				i.broke = true
				return
			}
		}
	}
}
