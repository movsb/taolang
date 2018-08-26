package main

import (
	"fmt"
)

type Statement interface {
	Execute(ctx *Context)
}

type Returner interface {
	Return() (*Value, bool)
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
	var value *Value
	if value = ctx.FindValue(v.Name, true); value == nil {
		panic(fmt.Sprintf("undefined variable: %s", v.Name))
	}
	*value = *v.Expr.Evaluate(ctx)
}

type FunctionStatement struct {
	name string
	expr *FunctionExpression
}

func (f *FunctionStatement) Execute(ctx *Context) {
	ctx.AddValue(f.name, ValueFromFunction(f.name, f.expr))
}

type ReturnStatement struct {
	expr  Expression
	value *Value
}

func (r *ReturnStatement) Execute(ctx *Context) {
	r.value = r.expr.Evaluate(ctx)
}

func (r *ReturnStatement) Return() (*Value, bool) {
	return r.value, true
}

type BlockStatement struct {
	retValue *Value
	broke    bool
	stmts    []Statement
}

func (b *BlockStatement) Return() (*Value, bool) {
	return b.retValue, b.retValue != nil
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

func (r *ExpressionStatement) Execute(ctx *Context) {
	value := r.expr.Evaluate(ctx)
	_ = value // drop expr value
}

type WhileStatement struct {
	expr     Expression
	block    *BlockStatement
	retValue *Value
}

func (w *WhileStatement) Return() (*Value, bool) {
	return w.retValue, w.retValue != nil
}

func (w *WhileStatement) Execute(ctx *Context) {
	for {
		if value := w.expr.Evaluate(ctx); !value.Truthy(ctx) {
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
