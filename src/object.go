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
	Functional(name string) *Builtin
}

// KeyAssigner is implemented by those who can be assigned.
type KeyAssigner interface {
	KeyAssign(key string, val Value)
}

// ElemAssigner is implemented by those who can be assigned.
type ElemAssigner interface {
	ElemAssign(elem int, val Value)
}

// Object is an object.
type Object struct {
	array *Array
	props map[string]Value
	ElemIndexer
	ElemAssigner
}

// NewObject news an object.
func NewObject() *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	o.ElemIndexer = nil
	o.ElemAssigner = nil
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

// KeyAssign implements KeyAssigner.
func (o *Object) KeyAssign(key string, val Value) {
	o.SetKey(key, val)
}

// IsArray tells if an object is acted as array.
func (o *Object) IsArray() bool {
	return o.ElemIndexer != nil
}

func (o *Object) String() string {
	if o.IsArray() {
		return o.array.String()
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
