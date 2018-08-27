package main

// Object is an object.
type Object struct {
	props map[string]Value
}

// NewObject news an object.
func NewObject() *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	return o
}

// Get gets a value by key.
func (o *Object) Get(key string) Value {
	if val, ok := o.props[key]; ok {
		return val
	}
	return Value{}
}

// Set sets a value by key.
func (o *Object) Set(key string, val Value) {
	o.props[key] = val
}
