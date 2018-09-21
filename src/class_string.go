package main

import (
	"strings"
)

// String is a string object.
type String struct {
	s string // the string
}

// NewString news a string object.
func NewString(s string) *String {
	return &String{s: s}
}

// GetKey implements KeyGetter.
func (s *String) GetKey(key string) Value {
	if fn, ok := _stringMethods[key]; ok {
		return ValueFromBuiltin(s, key, fn)
	}
	return ValueFromNil()
}

var _stringMethods map[string]BuiltinFunction

func init() {
	_stringMethods = map[string]BuiltinFunction{
		"lower": _stringLower,
	}
}

// Lower transforms the string into lower case.
func _stringLower(this interface{}, ctx *Context, args *Values) Value {
	return ValueFromString(strings.ToLower(this.(*String).s))
}
