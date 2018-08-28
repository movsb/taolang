package main

type KeyIndexer interface {
	Key(key string) *Value
	SetKey(key string, val Value)
}

type ElemIndexer interface {
	Len() int
	Elem(pos int) *Value
	SetElem(pos int, val Value)
	PushElem(val Value)
}

// Object is an object.
type Object struct {
	array interface{}
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
func NewArray() *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	o.array = &ValueArray{}
	o.ElemIndexer = (o.array).(ElemIndexer)
	return o
}

// Key gets a value by key.
func (o *Object) Key(key string) *Value {
	if o.IsArray() {
		if key == "length" {
			return ValueFromNumber(o.Len())
		}
	}
	if val, ok := o.props[key]; ok {
		return &val
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

type ValueArray struct {
	elems []Value
}

func (v *ValueArray) Len() int {
	return len(v.elems)
}

func (v *ValueArray) Elem(pos int) *Value {
	if pos < 0 || pos > len(v.elems)-1 {
		panic("array index out of range")
	}
	return &v.elems[pos]
}

func (v *ValueArray) SetElem(pos int, val Value) {
	if pos < 0 || pos > len(v.elems)-1 {
		panic("array index out of range")
	}
	v.elems[pos] = val
}

func (v *ValueArray) PushElem(val Value) {
	v.elems = append(v.elems, val)
}
