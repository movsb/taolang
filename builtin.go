package main

import (
	"fmt"
)

type Builtin func(ctx *Context, args Values) *Value

func InitBuiltins(ctx *Context) {
	pairs := [...]struct {
		name    string
		builtin Builtin
	}{
		{"print", print},
	}
	for _, pair := range pairs {
		ctx.AddValue(pair.name, ValueFromBuiltin(pair.name, pair.builtin))
	}
}

func print(ctx *Context, args Values) *Value {
	fmt.Print(args.ToInterfaces()...)
	return ValueFromNil()
}
