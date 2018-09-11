package main

import (
	"fmt"
	"time"
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
		"print":      ValueFromBuiltin(g, "print", _globalPrint),
		"println":    ValueFromBuiltin(g, "println", _globalPrintln),
		"setTimeout": ValueFromBuiltin(g, "setTimeout", _globalSetTimeout),
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

func _globalSetTimeout(this interface{}, ctx *Context, args *Values) Value {
	if args.Len() < 2 {
		panic("setTimeout: callback expected")
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

// Timer is a timer.
type Timer struct {
	timer *time.Timer
}

// NewTimer news a timer.
func NewTimer(ctx *Context, callback Value, timeout int) *Timer {
	t := time.NewTimer(time.Millisecond * time.Duration(timeout))
	go func() {
		select {
		case <-t.C:
			Sync(func() {
				CallFunc(ctx, callback)
			})
		}
	}()
	return &Timer{timer: t}
}

// Key implements KeyIndexer.
func (t *Timer) Key(key string) Value {
	if fn, ok := _timerMethods[key]; ok {
		return ValueFromBuiltin(t, key, fn)
	}
	return ValueFromNil()
}

// SetKey implements KeyIndexer.
func (t *Timer) SetKey(key string, val Value) {
	panic("not assignable")
}

var _timerMethods map[string]BuiltinFunction

func init() {
	_timerMethods = map[string]BuiltinFunction{
		"stop": _timerStop,
	}
}

func _timerStop(this interface{}, ctx *Context, args *Values) Value {
	t := this.(*Timer).timer
	return ValueFromBoolean(t.Stop())
}
