package main

// Promise is a promise.
type Promise struct {
	resolvedFunc  Value
	rejectedFunc  Value
	resolvedValue *Value
	rejectedValue *Value
	thenPromise   *Promise // if this promise is thened
	toPromise     *Promise // if the resolver/rejector returns a promise, forward to this
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
func (p *Promise) Resolve(resolvedValue Value) Value {
	p.resolvedValue = &resolvedValue
	Async(func() { p.invokeResolver() })
	return ValueFromNil()
}

// Reject rejects the promise.
func (p *Promise) Reject(rejectedValue Value) Value {
	p.rejectedValue = &rejectedValue
	Async(func() { p.invokeRejecter() })
	return ValueFromNil()
}

// Then chains promises.
func (p *Promise) Then(resolve Value, reject Value) Value {
	p.resolvedFunc = resolve
	p.rejectedFunc = reject
	np := &Promise{}
	p.thenPromise = np
	return ValueFromObject(np)
}

func (p *Promise) invokeResolver() {
	// forward to p.then
	if p.toPromise != nil {
		p.toPromise.resolvedValue = p.resolvedValue
		p.toPromise.invokeResolver()
		return
	}

	// not then-ed
	if p.resolvedFunc.isNil() {
		return
	}

	result := CallFunc(
		NewContext("--promise-resolve--", nil),
		p.resolvedFunc, p.resolvedValue,
	)

	if promise, ok := result.value.(*Promise); ok {
		promise.toPromise = p.thenPromise
	} else {
		p.thenPromise.Resolve(result)
	}
}

func (p *Promise) invokeRejecter() {
	// forward to p.then
	if p.toPromise != nil {
		p.toPromise.rejectedValue = p.rejectedValue
		p.toPromise.invokeRejecter()
		return
	}

	// not then-ed
	if p.rejectedFunc.isNil() {
		return
	}

	result := CallFunc(
		NewContext("--promise-reject--", nil),
		p.rejectedFunc, p.rejectedValue,
	)

	if promise, ok := result.value.(*Promise); ok {
		promise.toPromise = p.thenPromise
	} else {
		p.thenPromise.Reject(result)
	}
}

var _promiseMethods map[string]BuiltinFunction

func init() {
	_promiseMethods = map[string]BuiltinFunction{
		"then": _promiseThen,
	}
}

func _promiseResolve(this interface{}, ctx *Context, args *Values) Value {
	return this.(*Promise).Resolve(args.Shift())
}

func _promiseReject(this interface{}, ctx *Context, args *Values) Value {
	return this.(*Promise).Reject(args.Shift())
}

func _promiseThen(this interface{}, ctx *Context, args *Values) Value {
	return this.(*Promise).Then(args.Shift(), args.Shift())
}
