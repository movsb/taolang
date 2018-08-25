package main

import (
	"fmt"
)

type Builtin func(ctx *Context, args []*Value) *Value

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

func print(ctx *Context, args []*Value) *Value {
	fmt.Print(args)
	return ValueFromNil()
}
