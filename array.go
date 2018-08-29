package main

type Array struct {
	elems []Value
}

func (a *Array) Len() int {
	return len(a.elems)
}

func (a *Array) Elem(pos int) Value {
	if pos < 0 || pos > len(a.elems)-1 {
		panic("array index out of range")
	}
	return a.elems[pos]
}

func (a *Array) SetElem(pos int, val Value) {
	if pos < 0 || pos > len(a.elems)-1 {
		panic("array index out of range")
	}
	a.elems[pos] = val
}

func (a *Array) PushElem(val Value) {
	a.elems = append(a.elems, val)
}

/// functional programming implementations below

func (a *Array) _Call(ctx *Context, lambda *FunctionExpression, args ...Value) Value {
	lambda.BindArguments(ctx, args...)
	switch data := lambda.Execute(ctx); data.Type {
	case vtVariable:
		return ctx.MustFind(data.variable(), true)
	case vtFunction:
		newCtx := NewContext(ctx)
		fn := data.function()
		fn.BindArguments(newCtx, args...)
		return fn.Execute(newCtx)
	case vtBuiltin:
		newCtx := NewContext(ctx)
		builtin := data.builtin()
		return builtin.fn(newCtx, NewValues(args...))
	default:
		return data
	}
}

func (a *Array) _Each(callback func(elem Value) bool) {
	for i, n, next := 0, a.Len(), true; i < n && next; i++ {
		next = callback(a.elems[i])
	}
}

// Each iterates over a list of elements, yielding each in turn to an iteratee function.
func (a *Array) Each(ctx *Context, args *Values) Value {
	a._Each(func(elem Value) bool {
		newCtx := NewContext(ctx)
		lambda := args.values[0].function()
		a._Call(newCtx, lambda, elem)
		return true
	})
	return ValueFromNil()
}
