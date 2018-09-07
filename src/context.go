package main

// var usedContexts map[string]uint
//
// func init() {
// 	usedContexts = make(map[string]uint)
// }

// Symbol is a name value in the context scopes.
type Symbol struct {
	Name  string
	Value Value
}

// Context chains named values in call frames.
type Context struct {
	// the name of the context, for debug purpose
	name    string
	parent  *Context
	symbols []*Symbol
	broke   bool  // a break statement has executed
	hasret  bool  // Is retval set?
	retval  Value // a return statement has executed
}

// NewContext news a context from parent.
// name: who this context is created for.
// parent: the parent scope or the parent closure chain.
func NewContext(name string, parent *Context) *Context {
	return &Context{
		name:   name,
		parent: parent,
	}
}

// FindValue finds a value from context frames.
func (c *Context) FindValue(name string, outer bool) (Value, bool) {
	for _, symbol := range c.symbols {
		if symbol.Name == name {
			return symbol.Value, true
		}
	}
	if outer && c.parent != nil {
		return c.parent.FindValue(name, true)
	}
	return Value{}, false
}

// MustFind must find a named value.
// Upon failure, it panics.
func (c *Context) MustFind(name string, outer bool) Value {
	value, ok := c.FindValue(name, outer)
	if !ok {
		panicf("name `%s' not defined", name)
	}
	return value
}

// AddValue adds a new value in current context.
func (c *Context) AddValue(name string, value Value) {
	if _, ok := c.FindValue(name, false); ok {
		panicf("name `%s' redefined", name)
	}
	c.symbols = append(c.symbols, &Symbol{
		Name:  name,
		Value: value,
	})
}

// SetValue sets value of an existed name.
func (c *Context) SetValue(name string, value Value) {
	for _, symbol := range c.symbols {
		if symbol.Name == name {
			symbol.Value = value
			return
		}
	}
	if c.parent != nil {
		c.parent.SetValue(name, value)
		return
	}
	panicf("name `%s' not defined", name)
}

// SetParent sets parent context.
func (c *Context) SetParent(parent *Context) {
	c.parent = parent
}

// SetReturn sets block return value.
func (c *Context) SetReturn(retval Value) {
	c.hasret = true
	c.retval = retval
}

// SetBreak sets break flag.
func (c *Context) SetBreak() {
	c.broke = true
}
