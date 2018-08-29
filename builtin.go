package main

import (
	"fmt"
)

type Builtin struct {
	name string
	fn   func(ctx *Context, args *Values) Value
}

func NewBuiltin(name string, fn func(*Context, *Values) Value) *Builtin {
	return &Builtin{
		name: name,
		fn:   fn,
	}
}

func InitBuiltins(ctx *Context) {
	builtins := []Builtin{
		{"print", print},
		{"println", println},
	}
	for _, builtin := range builtins {
		ctx.AddValue(builtin.name, ValueFromBuiltin(&builtin))
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
