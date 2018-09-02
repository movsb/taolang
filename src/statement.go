package main

type Statement interface {
	Execute(ctx *Context)
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

// AssignmentStatement assigns right to left
// address(left) <- evaluate(right)
type AssignmentStatement struct {
	left  Expression
	right Expression
}

func (v *AssignmentStatement) Execute(ctx *Context) {
	addresser, ok := v.left.(Addresser)
	if !ok {
		val := v.left.Evaluate(ctx)
		panicf("not assignable: %v (type: %s)", val, val.TypeName())
	}
	ref := addresser.Address(ctx)
	if ref == nil {
		panic("cannot address")
	}
	*ref = v.right.Evaluate(ctx)
}

type FunctionStatement struct {
	expr *FunctionExpression
}

func (f *FunctionStatement) Execute(ctx *Context) {
	_ = f.expr.Evaluate(ctx)
}

type ReturnStatement struct {
	expr Expression
}

func NewReturnStatement(expr Expression) *ReturnStatement {
	return &ReturnStatement{
		expr: expr,
	}
}

func (r *ReturnStatement) Execute(ctx *Context) {
	retval := r.expr.Evaluate(ctx)
	ctx.SetReturn(retval)
}

type BlockStatement struct {
	stmts []Statement
}

func NewBlockStatement(stmts ...Statement) *BlockStatement {
	b := &BlockStatement{}
	for _, stmt := range stmts {
		b.stmts = append(b.stmts, stmt)
	}
	return b
}

func (b *BlockStatement) Execute(ctx *Context) {
	for _, stmt := range b.stmts {
		var newCtx *Context
		switch typed := stmt.(type) {
		case *BlockStatement:
			newCtx = NewContext("--block--", ctx)
			typed.Execute(newCtx)
		default:
			newCtx = ctx
			typed.Execute(ctx)
		}
		if newCtx.broke {
			ctx.broke = true
			break
		}
		if newCtx.hasret {
			ctx.hasret = true
			ctx.retval = newCtx.retval
			return
		}
	}
}

type ExpressionStatement struct {
	expr Expression
}

func (r *ExpressionStatement) Execute(ctx *Context) {
	value := r.expr.Evaluate(ctx)
	_ = value // drop expr value
}

// ForStatement simulates go-style for loop.
//
// for init; test; incr {
//      block
// }
type ForStatement struct {
	init  Statement
	test  Expression
	incr  interface{} // can be either Expression or Statement(without semicolon)
	block *BlockStatement
}

// Execute implements Statement.
func (f *ForStatement) Execute(ctx *Context) {
	if f.init != nil {
		f.init.Execute(ctx)
	}
	for {
		// test
		if f.test != nil {
			if !f.test.Evaluate(ctx).Truth(ctx) {
				break
			}
		}
		// block
		f.block.Execute(ctx)
		if ctx.hasret {
			return
		}
		if ctx.broke {
			break
		}
		// incr
		if f.incr != nil {
			if expr, ok := f.incr.(Expression); ok {
				expr.Evaluate(ctx)
			} else if stmt, ok := f.incr.(Statement); ok {
				stmt.Execute(ctx)
			}
		}
	}
}

type BreakStatement struct {
}

func (b *BreakStatement) Execute(ctx *Context) {
	ctx.SetBreak()
}

type IfStatement struct {
	cond      Expression
	ifBlock   *BlockStatement
	elseBlock Statement // if or block
}

func (i *IfStatement) Execute(ctx *Context) {
	cond := i.cond.Evaluate(ctx)
	if cond.Truth(ctx) {
		newCtx := NewContext("--block--", ctx)
		i.ifBlock.Execute(newCtx)
		if newCtx.broke {
			ctx.broke = true
			return
		}
		if newCtx.hasret {
			ctx.hasret = true
			ctx.retval = newCtx.retval
			return
		}
	} else {
		newCtx := NewContext("--block--", ctx)
		switch typed := i.elseBlock.(type) {
		case nil:
			return
		case *IfStatement:
			typed.Execute(newCtx)
		case *BlockStatement:
			typed.Execute(newCtx)
		default:
			panic("bad else stmt")
		}
		if newCtx.broke {
			ctx.broke = true
			return
		}
		if newCtx.hasret {
			ctx.hasret = true
			ctx.retval = newCtx.retval
			return
		}
	}
}
