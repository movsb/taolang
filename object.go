package main

type KeyIndexer interface {
	Key(key string) Value
	SetKey(key string, val Value)
}

type ElemIndexer interface {
	Len() int
	Elem(pos int) Value
	SetElem(pos int, val Value)
	PushElem(val Value)

	Each(ctx *Context, args *Values) Value
	Map(ctx *Context, args *Values) Value
	Reduce(ctx *Context, args *Values) Value
}

// Object is an object.
type Object struct {
	array *Array
	props map[string]Value
	ElemIndexer
}

// NewObject news an object.
func NewObject() *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	o.ElemIndexer = nil
	return o
}

// NewArray news an array.
func NewArray(elems ...Value) *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	o.array = &Array{elems: elems}
	o.array.object = o
	o.ElemIndexer = o.array
	return o
}

// Key gets a value by key.
func (o *Object) Key(key string) Value {
	if o.IsArray() {
		if key == "length" {
			return ValueFromNumber(o.Len())
		} else if key == "each" {
			return ValueFromBuiltin("each", o.Each)
		} else if key == "map" {
			return ValueFromBuiltin("map", o.Map)
		} else if key == "reduce" {
			return ValueFromBuiltin("reduce", o.Reduce)
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

func (o *Object) IsArray() bool {
	return o.ElemIndexer != nil
}
