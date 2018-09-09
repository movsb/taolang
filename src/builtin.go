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

// NewBuiltin news a Builtin.
func NewBuiltin(this interface{}, name string, fn BuiltinFunction) *Builtin {
	return &Builtin{
		this: this,
		name: name,
		fn:   fn,
	}
}

// Execute executes the builtin.
// This is not a Statement implementation.
func (b *Builtin) Execute(ctx *Context, args *Values) Value {
	return b.fn(b.this, ctx, args)
}

// NewGlobal news the global object.
func NewGlobal() *Global {
	g := &Global{}
	g.props = map[string]Value{
		"print":   ValueFromBuiltin(g, "print", _globalPrint),
		"println": ValueFromBuiltin(g, "println", _globalPrintln),
	}
	return g
}

// Global is the global object.
type Global struct {
	Object
}

// Key implements KeyIndexer.
func (g *Global) Key(key string) Value {
	if prop, ok := g.props[key]; ok {
		return prop
	}
	return ValueFromNil()
}

func _globalPrint(this interface{}, ctx *Context, args *Values) Value {
	fmt.Print(args.All()...)
	return ValueFromNil()
}

func _globalPrintln(this interface{}, ctx *Context, args *Values) Value {
	_globalPrint(this, ctx, args)
	_globalPrint(this, ctx, NewValues(ValueFromString("\n")))
	return ValueFromNil()
}
