package main

import (
	"fmt"
)

// BuiltinFunction is a language-supplied function.
type BuiltinFunction func(this interface{}, ctx *Context, args *Values) Value

// Builtin is a builtin function.
type Builtin struct {
	this interface{}     // The owner, if one exists
	name string          // The name of the builtin
	fn   BuiltinFunction // The function
}

// Execute executes the builtin.
// This is not a Statement implementation.
func (b *Builtin) Execute(ctx *Context, args *Values) Value {
	return b.fn(b.this, ctx, args)
}

// NewBuiltin news a Builtin.
func NewBuiltin(this interface{}, name string, fn func(interface{}, *Context, *Values) Value) *Builtin {
	return &Builtin{
		this: this,
		name: name,
		fn:   fn,
	}
}

// InitBuiltins registers builtins into ctx.
func InitBuiltins(ctx *Context) {
	builtins := []Builtin{
		{nil, "print", print},
		{nil, "println", println},
	}
	for _, b := range builtins {
		ctx.AddSymbol(b.name, ValueFromBuiltin(b.this, b.name, b.fn))
	}
}

func print(this interface{}, ctx *Context, args *Values) Value {
	fmt.Print(args.All()...)
	return ValueFromNil()
}

func println(this interface{}, ctx *Context, args *Values) Value {
	print(this, ctx, args)
	print(this, ctx, NewValues(ValueFromString("\n")))
	return ValueFromNil()
}
