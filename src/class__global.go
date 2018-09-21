package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Global is the global object.
type Global struct {
	Object
}

// NewGlobal news the global object.
func NewGlobal() *Global {
	g := &Global{}
	g.props = map[string]Value{
		"print":      ValueFromBuiltin(g, "print", _globalPrint),
		"println":    ValueFromBuiltin(g, "println", _globalPrintln),
		"setTimeout": ValueFromBuiltin(g, "setTimeout", _globalSetTimeout),
		"newPromise": ValueFromBuiltin(g, "newPromise", _globalNewPromise),
		"newChannel": ValueFromBuiltin(g, "newChannel", _globalNewChannel),
		"httpGet":    ValueFromBuiltin(g, "httpGet", _globalHTTPGet),
	}
	return g
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
	args.Push(ValueFromString("\n"))
	_globalPrint(this, ctx, args)
	return ValueFromNil()
}

func _globalSetTimeout(this interface{}, ctx *Context, args *Values) Value {
	if args.Len() < 1 {
		panic(NewTypeError("setTimeout: callback expected"))
	} else if args.Len() < 2 {
		panic(NewTypeError("setTimeout: timeout expected"))
	}
	var callback = args.Shift()
	if !callback.isCallable() {
		panic(NewTypeError("setTimeout: callback must be a callable"))
	}
	timeout := args.Shift()
	if !timeout.isNumber() {
		panic(NewTypeError("setTimeout: timeout must be a number"))
	}
	t := NewTimer(ctx, callback, timeout.number())
	return ValueFromObject(t)
}

func _globalNewPromise(this interface{}, ctx *Context, args *Values) Value {
	return ValueFromObject(NewPromise(args.At(0)))
}

func _globalNewChannel(this interface{}, ctx *Context, args *Values) Value {
	return ValueFromObject(NewChannel(args.At(0)))
}

func _globalHTTPGet(this interface{}, ctx *Context, args *Values) Value {
	promise := &Promise{}
	go func() {
		url := args.Shift().str()
		resp, err := http.Get(url)
		if err != nil {
			// TODO
			panic(err)
		}

		bys, _ := ioutil.ReadAll(resp.Body)
		Async(func() {
			promise.Resolve(ValueFromString(string(bys)))
		})
	}()
	return ValueFromObject(promise)
}
