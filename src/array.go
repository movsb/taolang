package main

type Array struct {
	object *Object
	elems  []Value
	funcs  map[string]func(*Context, *Values) Value
}

// NewArray news an array.
func NewArray(elems ...Value) *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	o.array = &Array{elems: elems}
	o.ElemIndexer = o.array
	o.array.object = o
	o.array.funcs = map[string]func(*Context, *Values) Value{
		"each":   o.array.Each,
		"map":    o.array.Map,
		"reduce": o.array.Reduce,
		"find":   o.array.Find,
		"filter": o.array.Filter,
		"where":  o.array.Where,
	}
	return o
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

func (a *Array) Functional(name string) *Builtin {
	if fn, ok := a.funcs[name]; ok {
		return &Builtin{name: name, fn: fn}
	}
	return nil
}

func (a *Array) _Call(ctx *Context, lambdaValue Value, args ...Value) Value {
	ctx = NewContext("--lambda--", nil)
	lambda := lambdaValue.function()
	lambda.BindArguments(ctx, args...)
	switch data := lambda.Execute(ctx); data.Type {
	case vtVariable:
		return ctx.MustFind(data.variable(), true)
	case vtFunction:
		fn := data.function()
		newCtx := NewContext(fn.expr.name, nil)
		fn.BindArguments(newCtx, args...)
		return fn.Execute(newCtx)
	case vtBuiltin:
		builtin := data.builtin()
		newCtx := NewContext(builtin.name, nil)
		return builtin.fn(newCtx, NewValues(args...))
	default:
		return data
	}
}

func (a *Array) _Each(callback func(elem Value, index Value) bool) {
	for i, n, next := 0, a.Len(), true; i < n && next; i++ {
		next = callback(a.elems[i], ValueFromNumber(i))
	}
}

// Each iterates over a list of elements, yielding each in turn to an iteratee function.
func (a *Array) Each(ctx *Context, args *Values) Value {
	object := ValueFromObject(a.object)
	a._Each(func(elem Value, index Value) bool {
		a._Call(ctx, args.At(0), elem, index, object)
		return true
	})
	return ValueFromNil()
}

// Map produces a new array of values by mapping each value.
func (a *Array) Map(ctx *Context, args *Values) Value {
	object := ValueFromObject(a.object)
	values := make([]Value, 0, a.Len())
	a._Each(func(elem Value, index Value) bool {
		data := a._Call(ctx, args.At(0), elem, index, object)
		values = append(values, data)
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// Reduce boils down the array into a single value.
func (a *Array) Reduce(ctx *Context, args *Values) Value {
	if args.Len() < 2 {
		panic("usage: reduce(lambda, init)")
	}
	object := ValueFromObject(a.object)
	memo := args.At(1)
	a._Each(func(elem Value, index Value) bool {
		memo = a._Call(ctx, args.At(0), memo, elem, index, object)
		return true
	})
	return memo
}

// Find finds the first value.
func (a *Array) Find(ctx *Context, args *Values) Value {
	found := Value{}
	a._Each(func(elem Value, index Value) bool {
		if a._Call(ctx, args.At(0), elem).Truth(ctx) {
			found = elem
			return false
		}
		return true
	})
	return found
}

// Filter filters values.
func (a *Array) Filter(ctx *Context, args *Values) Value {
	values := make([]Value, 0, a.Len())
	a._Each(func(elem Value, index Value) bool {
		if a._Call(ctx, args.At(0), elem).Truth(ctx) {
			values = append(values, elem)
		}
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// Where filters objects by column conditions.
// save as Filter currently.
func (a *Array) Where(ctx *Context, args *Values) Value {
	values := make([]Value, 0, a.Len())
	a._Each(func(elem Value, index Value) bool {
		if a._Call(ctx, args.At(0), elem).Truth(ctx) {
			values = append(values, elem)
		}
		return true
	})
	return ValueFromObject(NewArray(values...))
}
