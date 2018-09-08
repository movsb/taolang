package main

import (
	"fmt"
)

// Builtin is a builtin function.
type Builtin struct {
	name string
	fn   func(ctx *Context, args *Values) Value
}

// NewBuiltin news a Builtin.
func NewBuiltin(name string, fn func(*Context, *Values) Value) *Builtin {
	return &Builtin{
		name: name,
		fn:   fn,
	}
}

// InitBuiltins registers builtins into ctx.
func InitBuiltins(ctx *Context) {
	builtins := []Builtin{
		{"print", print},
		{"println", println},
	}
	for _, b := range builtins {
		ctx.AddSymbol(b.name, ValueFromBuiltin(b.name, b.fn))
	}
}

func print(ctx *Context, args *Values) Value {
	fmt.Print(args.All()...)
	return ValueFromNil()
}

func println(ctx *Context, args *Values) Value {
	print(ctx, args)
	print(ctx, NewValues(ValueFromString("\n")))
	return ValueFromNil()
}
