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
	Functional(name string) *Builtin
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

// Key gets a value by key.
func (o *Object) Key(key string) Value {
	if o.IsArray() {
		if key == "length" {
			return ValueFromNumber(o.Len())
		}
		if builtin := o.Functional(key); builtin != nil {
			return ValueFromBuiltin(builtin.name, builtin.fn)
		}
		return Value{}
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
