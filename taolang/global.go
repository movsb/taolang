package taolang

import (
	"fmt"
	"io"
)

// Global is the global object.
type Global struct {
	props map[string]Value
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

// GetProp implements IObject.
func (g *Global) GetProp(key string) Value {
	if prop, ok := g.props[key]; ok {
		return prop
	}
	return ValueFromNil()
}

// SetProp implements IObject.
func (g *Global) SetProp(key string, val Value) {
	g.props[key] = val
}

var Stdout io.Writer

func _globalPrint(this interface{}, ctx *Context, args *Values) Value {
	if Stdout != nil {
		fmt.Fprint(Stdout, args.All()...)
	} else {
		fmt.Print(args.All()...)
	}
	return ValueFromNil()
}

func _globalPrintln(this interface{}, ctx *Context, args *Values) Value {
	if Stdout != nil {
		fmt.Fprintln(Stdout, args.All()...)
	} else {
		fmt.Println(args.All()...)
	}
	return ValueFromNil()
}
