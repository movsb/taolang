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
		"newPromise": ValueFromBuiltin(g, "newPromise", _globalNewPromise),
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
	return ValueFromObject(NewPromise(ctx, args.At(0)))
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
			Async(func() {
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

// Promise is a promise.
type Promise struct {
	resolvedValue *Value
	rejectedValue *Value
	resolvedFunc  Value
	rejectedFunc  Value
	thenPromise   *Promise
}

// NewPromise news a promise.
func NewPromise(ctx *Context, executor Value) *Promise {
	promise := &Promise{}
	resolve := ValueFromBuiltin(promise, "resolve", _promiseResolve)
	reject := ValueFromBuiltin(promise, "reject", _promiseReject)
	CallFunc(NewContext("--promise-executor--", ctx), executor, resolve, reject)
	return promise
}

// Key implements KeyIndexer.
func (p *Promise) Key(key string) Value {
	if fn, ok := _promiseMethods[key]; ok {
		return ValueFromBuiltin(p, key, fn)
	}
	return ValueFromNil()
}

// SetKey implements KeyIndexer.
func (p *Promise) SetKey(key string, val Value) {
	panic("not assignable")
}

// Resolve resolves the promise.
func (p *Promise) Resolve(ctx *Context, args *Values) Value {
	v := args.Shift() // maybe no value
	p.resolvedValue = &v
	Async(func() {
		if p.resolvedFunc.isFunction() || p.resolvedFunc.isBuiltin() {
			result := CallFunc(
				NewContext("--promise-resolve--", ctx),
				p.resolvedFunc, p.resolvedValue,
			)
			if p.thenPromise != nil {
				if promise, ok := result.value.(*Promise); ok {
					if promise.resolvedValue == nil && promise.rejectedValue == nil {
						promise.resolvedFunc = p.thenPromise.resolvedFunc
						promise.rejectedFunc = p.thenPromise.rejectedFunc
					} else {
						if promise.resolvedValue != nil {
							p.thenPromise.Resolve(ctx, NewValues(*promise.resolvedValue))
						}
						if promise.rejectedValue != nil {
							p.thenPromise.Reject(ctx, NewValues(*promise.rejectedValue))
						}
					}
				} else {
					p.thenPromise.Resolve(ctx, NewValues(result))
				}
			}
		}
	})
	return ValueFromNil()
}

// Reject rejects the promise.
func (p *Promise) Reject(ctx *Context, args *Values) Value {
	v := args.Shift() // maybe no value
	p.rejectedValue = &v
	Async(func() {
		if p.rejectedFunc.isFunction() || p.rejectedFunc.isBuiltin() {
			result := CallFunc(
				NewContext("--promise-reject--", ctx),
				p.rejectedFunc, p.rejectedValue,
			)
			if p.thenPromise != nil {
				if promise, ok := result.value.(*Promise); ok {
					if promise.resolvedValue == nil && promise.rejectedValue == nil {
						promise.resolvedFunc = p.thenPromise.resolvedFunc
						promise.rejectedFunc = p.thenPromise.rejectedFunc
					} else {
						if promise.resolvedValue != nil {
							p.thenPromise.Resolve(ctx, NewValues(*promise.resolvedValue))
						}
						if promise.rejectedValue != nil {
							p.thenPromise.Reject(ctx, NewValues(*promise.rejectedValue))
						}
					}
				} else {
					p.thenPromise.Reject(ctx, NewValues(result))
				}
			}
		}
	})
	return ValueFromNil()
}

// Then chains promises.
func (p *Promise) Then(ctx *Context, args *Values) Value {
	if args.Len() >= 1 {
		p.resolvedFunc = args.Shift()
	}
	if args.Len() >= 1 {
		p.rejectedFunc = args.Shift()
	}
	np := &Promise{}
	p.thenPromise = np
	return ValueFromObject(np)
}

var _promiseMethods map[string]BuiltinFunction

func init() {
	_promiseMethods = map[string]BuiltinFunction{
		"then": _promiseThen,
	}
}

func _promiseResolve(this interface{}, ctx *Context, args *Values) Value {
	return this.(*Promise).Resolve(ctx, args)
}

func _promiseReject(this interface{}, ctx *Context, args *Values) Value {
	return this.(*Promise).Reject(ctx, args)
}

func _promiseThen(this interface{}, ctx *Context, args *Values) Value {
	return this.(*Promise).Then(ctx, args)
}
