package main

// Statement is implemented by those who can execute.
type Statement interface {
	Execute(ctx *Context)
}

// VariableStatement is a let ... statement, which defines and inits variables.
type VariableStatement struct {
	Name string
	Expr Expression
}

// Execute implements Statement.
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

// Execute implements Statement.
func (v *AssignmentStatement) Execute(ctx *Context) {
	assigner, ok := v.left.(Assigner)
	if !ok {
		val := v.left.Evaluate(ctx)
		panicf("not assignable: %v (type: %s)", val, val.TypeName())
	}
	value := v.right.Evaluate(ctx)
	assigner.Assign(ctx, value)
}

// FunctionStatement is a function definition statement.
// Because Tao treats functions as first-class values.
// We call it a function expression.
type FunctionStatement struct {
	expr *FunctionExpression
}

// Execute implements Statement.
func (f *FunctionStatement) Execute(ctx *Context) {
	_ = f.expr.Evaluate(ctx)
}

// ReturnStatement is the `return ...;` statement.
type ReturnStatement struct {
	expr Expression
}

// NewReturnStatement news a ReturnStatement.
func NewReturnStatement(expr Expression) *ReturnStatement {
	return &ReturnStatement{
		expr: expr,
	}
}

// Execute implements Statement.
func (r *ReturnStatement) Execute(ctx *Context) {
	retval := ValueFromNil()
	if r.expr != nil {
		retval = r.expr.Evaluate(ctx)
	}
	ctx.SetReturn(retval)
}

// BlockStatement is a `{ ... }` statement.
type BlockStatement struct {
	stmts []Statement
}

// NewBlockStatement news a block statement.
func NewBlockStatement(stmts ...Statement) *BlockStatement {
	b := &BlockStatement{}
	for _, stmt := range stmts {
		b.stmts = append(b.stmts, stmt)
	}
	return b
}

// Execute implements Statement.
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

// ExpressionStatement is the expression statement.
// An expression statement is an expression ended with a semicolon.
// The evaluated value of the expression is simply dropped.
type ExpressionStatement struct {
	expr Expression
}

// Execute implements Statement.
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
		newCtx := NewContext("--for-block--", ctx)
		f.block.Execute(newCtx)
		if newCtx.hasret {
			ctx.SetReturn(newCtx.retval)
			return
		}
		if newCtx.broke {
			ctx.SetBreak()
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

// BreakStatement is the break statement.
type BreakStatement struct {
}

// Execute implements Statement.
func (b *BreakStatement) Execute(ctx *Context) {
	ctx.SetBreak()
}

// IfStatement is the if ... else if ... statement.
type IfStatement struct {
	cond      Expression
	ifBlock   *BlockStatement
	elseBlock Statement // if or block
}

// Execute implements Statement.
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
