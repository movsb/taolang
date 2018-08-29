package main

// Symbol is a name value in the context scopes.
type Symbol struct {
	Name  string
	Value Value
}

// Context chains named values in call frames.
type Context struct {
	parent  *Context
	symbols []*Symbol
}

// NewContext news a context from parent.
func NewContext(parent *Context) *Context {
	return &Context{
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
		panicf("name(%s) not found", name)
	}
	return value
}

// AddValue adds a new value in current context.
func (c *Context) AddValue(name string, value Value) {
	if _, ok := c.FindValue(name, false); ok {
		panicf("name `%s' is already defined in this scope", name)
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
	panicf("name `%s' is not defined", name)
}
