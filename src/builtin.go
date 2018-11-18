package main

// BuiltinFunction is a language-supplied function.
type BuiltinFunction func(this interface{}, ctx *Context, args *Values) Value

// Builtin is a builtin function.
type Builtin struct {
	this interface{}     // The owner, if one exists
	name string          // The name of the builtin
	fn   BuiltinFunction // The function
}

// NewBuiltin news a Builtin.
func NewBuiltin(this interface{}, name string, fn BuiltinFunction) *Builtin {
	return &Builtin{
		this: this,
		name: name,
		fn:   fn,
	}
}

// Execute executes the builtin.
// It implements Callable.
func (b *Builtin) Execute(ctx *Context, args *Values) Value {
	return b.fn(b.this, ctx, args)
}

// Constructable is constructable.
type Constructable interface {
	Construct(ctx *Context, args *Values) IObject
}

// Constructor is a class constructor.
type Constructor struct {
	Name string
	Ctor Constructable
}

// BuiltinConstructor is a builtin constructor.
type BuiltinConstructor struct {
	ctor BuiltinFunction
}

// Construct implements Constructable.
func (c BuiltinConstructor) Construct(ctx *Context, args *Values) IObject {
	return c.ctor(nil, ctx, args).object()
}
