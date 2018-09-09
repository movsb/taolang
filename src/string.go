package main

import (
	"strings"
)

// String is a string object.
type String struct {
	s string
}

// NewString news a string object.
func NewString(s string) *String {
	return &String{s: s}
}

// Key implements KeyIndexer.
func (s *String) Key(key string) Value {
	if fn, ok := _globalStringMethods[key]; ok {
		return ValueFromBuiltin(s, key, fn)
	}
	return ValueFromNil()
}

// SetKey implements KeyIndexer.
func (s *String) SetKey(key string, val Value) {
	panic("not assignable")
}

var _globalStringMethods map[string]BuiltinFunction

func init() {
	_globalStringMethods = map[string]BuiltinFunction{
		"lower": _stringLower,
	}
}

// Lower transforms the string into lower case.
func _stringLower(this interface{}, ctx *Context, args *Values) Value {
	return ValueFromString(strings.ToLower(this.(*String).s))
}
