package main

import (
	"strings"
)

// String is a string object.
type String struct {
	s string // the string
	c []rune // the characters
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

// Len implements ElemGetter & ElemSetter.
func (s *String) Len() int {
	s.initChars()
	return len(s.c)
}

// GetElem implements ElemGetter.
func (s *String) GetElem(pos int) Value {
	s.initChars()
	if pos < 0 || pos > len(s.c)-1 {
		panic(NewRangeError("character index out of range"))
	}
	return ValueFromString(string(s.c[pos]))
}

func (s *String) initChars() {
	if s.c == nil {
		s.c = []rune(s.s)
	}
}

var _stringMethods map[string]BuiltinFunction

func init() {
	_stringMethods = map[string]BuiltinFunction{
		"len":   _stringLen,
		"lower": _stringLower,
	}
}

// Len lens the string.
func _stringLen(this interface{}, ctx *Context, args *Values) Value {
	return ValueFromNumber(this.(*String).Len())
}

// Lower transforms the string into lower case.
func _stringLower(this interface{}, ctx *Context, args *Values) Value {
	return ValueFromString(strings.ToLower(this.(*String).s))
}
