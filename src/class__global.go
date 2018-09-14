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
	_globalPrint(this, ctx, args)
	_globalPrint(this, ctx, NewValues(ValueFromString("\n")))
	return ValueFromNil()
}

func _globalSetTimeout(this interface{}, ctx *Context, args *Values) Value {
	if args.Len() < 1 {
		panic("setTimeout: callback expected")
	} else if args.Len() < 2 {
		panic("setTimeout: timeout expected")
	}
	var callback = args.Shift()
	if !callback.isFunction() && !callback.isBuiltin() {
		panic("setTimeout: callback must be a function")
	}
	timeout := args.Shift()
	if !timeout.isNumber() {
		panic("setTimeout: timeout must be a number")
	}
	t := NewTimer(ctx, callback, timeout.number())
	return ValueFromObject(t)
}

func _globalNewPromise(this interface{}, ctx *Context, args *Values) Value {
	return ValueFromObject(NewPromise(args.At(0)))
}

func _globalHTTPGet(this interface{}, ctx *Context, args *Values) Value {
	promise := &Promise{}
	go func() {
		url := args.Shift().str()
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		bys, _ := ioutil.ReadAll(resp.Body)
		Async(func() {
			promise.Resolve(ValueFromString(string(bys)))
		})
	}()
	return ValueFromObject(promise)
}
