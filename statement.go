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
	expr *FunctionDefinitionExpression
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

type ExpressionStatement struct {
	expr Expression
}

func (r *ExpressionStatement) Execute(ctx *Context) {
	value := r.expr.Evaluate(ctx)
	_ = value // drop expr value
}
