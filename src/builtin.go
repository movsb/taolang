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
// This is not a Statement implementation.
func (b *Builtin) Execute(ctx *Context, args *Values) Value {
	return b.fn(b.this, ctx, args)
}
