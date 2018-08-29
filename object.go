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
func NewArray() *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	o.array = &Array{}
	o.ElemIndexer = o.array
	return o
}

// Key gets a value by key.
func (o *Object) Key(key string) Value {
	if o.IsArray() {
		if key == "length" {
			return ValueFromNumber(o.Len())
		} else if key == "each" {
			return ValueFromBuiltin(NewBuiltin("each", o.Each))
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

func (a *Array) Each(ctx *Context, args *Values) Value {
	if args.Len() != 1 || !args.values[0].isFunction() {
		panic("each accepts function or lambda only")
	}
	fn := args.values[0].function()
	if fn.params.Len() != 1 {
		panic("one parameter only")
	}
	for i, n := 0, a.Len(); i < n; i++ {
		newCtx := NewContext(ctx)
		newCtx.AddValue(fn.params.names[0], a.elems[i])
		maybeVar := fn.Execute(newCtx)
		switch maybeVar.Type {
		case vtVariable:
			value, ok := newCtx.FindValue(maybeVar.variable(), true)
			if !ok {
				panic("variable not defined")
			}
			_ = value // drop it because each doesn't need it
		case vtFunction:
			f := maybeVar.function()
			if f.params.Len() != 1 {
				panic("one parameter only")
			}
			newCtx2 := NewContext(newCtx)
			newCtx2.AddValue(f.params.names[0], a.elems[i])
			_ = f.Execute(newCtx2)
		case vtBuiltin:
			newCtx2 := NewContext(newCtx)
			values := NewValues(a.elems[i])
			_ = maybeVar.builtin().fn(newCtx2, values)
		}

	}
	return ValueFromNil()
}
