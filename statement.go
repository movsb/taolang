package main

import (
	"fmt"
)

type Statement interface {
	Execute(ctx *Context)
}

type VariableDefinitionStatement struct {
	Name string
	Expr Expression
}

func (v *VariableDefinitionStatement) Execute(ctx *Context) {
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

type FunctionDefinitionStatement struct {
	name string
	expr *FunctionExpression
}

func (f *FunctionDefinitionStatement) Execute(ctx *Context) {
	ctx.AddValue(f.name, ValueFromFunction(f.name, f.expr))
}

type ReturnStatement struct {
	expr  Expression
	value *Value
}

func (r *ReturnStatement) Execute(ctx *Context) {
	r.value = r.expr.Evaluate(ctx)
}

type BlockStatement struct {
	retValue *Value
	stmts    []Statement
}

func (b *BlockStatement) Returned() (value *Value, hasReturned bool) {
	return b.retValue, b.retValue != nil
}

func (b *BlockStatement) Execute(ctx *Context) {
	for _, stmt := range b.stmts {
		switch typed := stmt.(type) {
		case *BlockStatement:
			newCtx := NewContext(ctx)
			typed.Execute(newCtx)
			if ret, ok := typed.Returned(); ok {
				b.retValue = ret
				return
			}
		case *ReturnStatement:
			typed.Execute(ctx)
			b.retValue = typed.value
			return
		default:
			typed.Execute(ctx)
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
