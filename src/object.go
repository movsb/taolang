package main

import (
	"bytes"
	"fmt"
)

// KeyIndexer is implemented by those who supports keyed values.
type KeyIndexer interface {
	Key(key string) Value
	SetKey(key string, val Value)
}

// ElemIndexer is implemented by those who supports indexed values.
type ElemIndexer interface {
	Len() int
	Elem(pos int) Value
	SetElem(pos int, val Value)
	PushElem(val Value)
}

// KeyAssigner is implemented by those who can be assigned.
type KeyAssigner interface {
	KeyAssign(key string, val Value)
}

// ElemAssigner is implemented by those who can be assigned.
type ElemAssigner interface {
	ElemAssign(elem int, val Value)
}

// Object is either an object or an array.
type Object struct {
	elems []Value          // array elements
	props map[string]Value // object properties
	array bool
}

// NewObject news an object.
func NewObject() *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	return o
}

// NewArray news an array.
func NewArray(elems ...Value) *Object {
	o := NewObject()
	o.array = true
	o.elems = elems

	builtins := map[string]func(*Context, *Values) Value{
		"each":   o.Each,
		"map":    o.Map,
		"reduce": o.Reduce,
		"find":   o.Find,
		"filter": o.Filter,
		"where":  o.Where,
	}

	for k, v := range builtins {
		o.props[k] = ValueFromBuiltin(k, v)
	}

	return o
}

// Key gets a value by key.
func (o *Object) Key(key string) Value {
	if o.array {
		if key == "length" {
			return ValueFromNumber(o.Len())
		}
	}
	if val, ok := o.props[key]; ok {
		return val
	}
	return ValueFromNil()
}

// SetKey sets a value by key.
func (o *Object) SetKey(key string, val Value) {
	o.props[key] = val
}

// KeyAssign implements KeyAssigner.
func (o *Object) KeyAssign(key string, val Value) {
	o.SetKey(key, val)
}

// Len implements ElemIndexer.
func (o *Object) Len() int {
	return len(o.elems)
}

// Elem implements ElemIndexer.
func (o *Object) Elem(pos int) Value {
	if pos < 0 || pos > len(o.elems)-1 {
		panic("array index out of range")
	}
	return o.elems[pos]
}

// SetElem implements ElemIndexer.
func (o *Object) SetElem(pos int, val Value) {
	if pos < 0 || pos > len(o.elems)-1 {
		panic("array index out of range")
	}
	o.elems[pos] = val
}

// ElemAssign implements ElemAssigner.
func (o *Object) ElemAssign(elem int, val Value) {
	o.SetElem(elem, val)
}

// PushElem implements ElemIndexer.
func (o *Object) PushElem(val Value) {
	o.elems = append(o.elems, val)
}

func (o *Object) String() string {
	if o.array {
		buf := bytes.NewBuffer(nil)
		buf.WriteString("[")
		for i, n := 0, o.Len(); i < n; i++ {
			elem := o.Elem(i)
			buf.WriteString(elem.String())
			if i != n-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("]")
		return buf.String()
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString("{")
	n := len(o.props)
	i := 0
	for k, p := range o.props {
		// TODO k may have invalid characters.
		buf.WriteString(fmt.Sprintf(`%s:%v`, k, p))
		if i != n-1 {
			buf.WriteString(",")
		}
		i++
	}
	buf.WriteString("}")
	return buf.String()
}

func (o *Object) _Call(ctx *Context, lambdaValue Value, args ...Value) Value {
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

func (o *Object) _Each(callback func(elem Value, index Value) bool) {
	for i, n, next := 0, o.Len(), true; i < n && next; i++ {
		next = callback(o.elems[i], ValueFromNumber(i))
	}
}

// Each iterates over a list of elements, yielding each in turn to an iteratee function.
func (o *Object) Each(ctx *Context, args *Values) Value {
	object := ValueFromObject(o)
	o._Each(func(elem Value, index Value) bool {
		o._Call(ctx, args.At(0), elem, index, object)
		return true
	})
	return ValueFromNil()
}

// Map produces a new array of values by mapping each value.
func (o *Object) Map(ctx *Context, args *Values) Value {
	object := ValueFromObject(o)
	values := make([]Value, 0, o.Len())
	o._Each(func(elem Value, index Value) bool {
		data := o._Call(ctx, args.At(0), elem, index, object)
		values = append(values, data)
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// Reduce boils down the array into a single value.
func (o *Object) Reduce(ctx *Context, args *Values) Value {
	if args.Len() < 2 {
		panic("usage: reduce(lambda, init)")
	}
	object := ValueFromObject(o)
	memo := args.At(1)
	o._Each(func(elem Value, index Value) bool {
		memo = o._Call(ctx, args.At(0), memo, elem, index, object)
		return true
	})
	return memo
}

// Find finds the first value.
func (o *Object) Find(ctx *Context, args *Values) Value {
	found := Value{}
	o._Each(func(elem Value, index Value) bool {
		if o._Call(ctx, args.At(0), elem).Truth(ctx) {
			found = elem
			return false
		}
		return true
	})
	return found
}

// Filter filters values.
func (o *Object) Filter(ctx *Context, args *Values) Value {
	values := make([]Value, 0, o.Len())
	o._Each(func(elem Value, index Value) bool {
		if o._Call(ctx, args.At(0), elem).Truth(ctx) {
			values = append(values, elem)
		}
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// Where filters objects by column conditions.
// save as Filter currently.
func (o *Object) Where(ctx *Context, args *Values) Value {
	values := make([]Value, 0, o.Len())
	o._Each(func(elem Value, index Value) bool {
		if o._Call(ctx, args.At(0), elem).Truth(ctx) {
			values = append(values, elem)
		}
		return true
	})
	return ValueFromObject(NewArray(values...))
}
