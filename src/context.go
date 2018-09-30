package main

// Symbol is a named value in the context.
type Symbol struct {
	Name  string // symbol name
	Value Value  // symbol value
}

// Context is a place where named symbols are saved for current scope.
// A context will be created when:
//   - calls a tao function (includes lambda)
//   - calls a builtin function
//   - executes a block statement (if, for, etc.)
type Context struct {
	name    string    // the name of the Context, for debug only
	parent  *Context  // the parent of the context, for scope chain
	symbols []*Symbol // symbols defined in this scope
	broke   bool      // a break statement has executed
	hasret  bool      // Is retval set?
	retval  Value     // a return statement has executed
}

// NewContext news a context.
// name: whom this context is created for.
// parent: the parent scope or the parent closure chain.
func NewContext(name string, parent *Context) *Context {
	return &Context{
		name:   name,
		parent: parent,
	}
}

// FindSymbol finds a symbol from context chain.
func (c *Context) FindSymbol(name string, outer bool) (Value, bool) {
	for _, symbol := range c.symbols {
		if symbol.Name == name {
			return symbol.Value, true
		}
	}
	// If not found, find outer scope.
	if outer {
		if c.parent != nil {
			return c.parent.FindSymbol(name, true)
		}
		return c.FromGlobal(name)
	}
	return Value{}, false
}

// MustFind must find a symbol.
// Upon failure, it panics.
func (c *Context) MustFind(name string, outer bool) Value {
	value, ok := c.FindSymbol(name, outer)
	if !ok {
		panic(NewNameError("name `%s' not defined", name))
	}
	return value
}

// FromGlobal finds a symbol from global.
func (c *Context) FromGlobal(name string) (Value, bool) {
	// This is the global context
	// TODO use global directly
	global := c.MustFind("global", false)
	if !global.isObject() {
		panic(NewTypeError("global is not an object"))
	}
	if obj, ok := global.object().(*Global); ok {
		val, ok := obj.Object.props[name]
		return val, ok
	}
	return Value{}, false
}

// AddSymbol adds a new symbol in current context.
// If a symbol with given name does exist, It will panic.
func (c *Context) AddSymbol(name string, value Value) {
	if _, ok := c.FindSymbol(name, false); ok {
		panic(NewNameError("name `%s' redefined", name))
	}
	c.symbols = append(c.symbols, &Symbol{
		Name:  name,
		Value: value,
	})
}

// AddObject adds an object into context.
func (c *Context) AddObject(name string, obj KeyIndexer) {
	c.AddSymbol(name, ValueFromObject(obj))
}

// AddClass adds a callable into context.
func (c *Context) AddClass(name string, ctor Constructable) {
	c.AddSymbol(name, ValueFromClass(name, ctor))
}

// SetSymbol sets the value of a symbol.
func (c *Context) SetSymbol(name string, value Value) {
	for _, symbol := range c.symbols {
		if symbol.Name == name {
			symbol.Value = value
			return
		}
	}
	if c.parent != nil {
		c.parent.SetSymbol(name, value)
		return
	}
	panic(NewNameError("name `%s' not defined", name))
}

// SetParent sets the parent context.
func (c *Context) SetParent(parent *Context) {
	c.parent = parent
}

// SetReturn sets the block return value.
func (c *Context) SetReturn(retval Value) {
	c.hasret = true
	c.retval = retval
}

// SetBreak sets the break flag.
func (c *Context) SetBreak() {
	c.broke = true
}
