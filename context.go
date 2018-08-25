package main

import (
	"fmt"
)

type Symbol struct {
	Name  string
	Value *Value
}

type Context struct {
	parent  *Context
	symbols []*Symbol
}

func NewContext(parent *Context) *Context {
	return &Context{
		parent: parent,
	}
}

func (c *Context) FindValue(name string, outer bool) *Value {
	for _, symbol := range c.symbols {
		if symbol.Name == name {
			return symbol.Value
		}
	}
	if outer && c.parent != nil {
		return c.parent.FindValue(name, true)
	}
	return nil
}

func (c *Context) AddValue(name string, value *Value) {
	if c.FindValue(name, false) != nil {
		panic(fmt.Sprintf("name `%s' is already defined in this scope", name))
	}
	c.symbols = append(c.symbols, &Symbol{
		Name:  name,
		Value: value,
	})
}

func (c *Context) SetValue(name string, value *Value) {
	exist := c.FindValue(name, true)
	if exist == nil {
		panic(fmt.Sprintf("name `%s' is not defined", name))
	}
	*exist = *value
}
