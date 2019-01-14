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

// GetProp implements IObject.
func (s *String) GetProp(key string) Value {
	if fn, ok := _stringMethods[key]; ok {
		return ValueFromBuiltin(s, key, fn)
	}
	return ValueFromNil()
}

// SetProp implements IObject.
func (s *String) SetProp(key string, val Value) {
	panic(NewNotAssignableError(ValueFromString(s.s)))
}

// Len implements IArray.
func (s *String) Len() int {
	s.initChars()
	return len(s.c)
}

// GetElem implements IArray.
func (s *String) GetElem(pos int) Value {
	s.initChars()
	if pos < 0 || pos > len(s.c)-1 {
		panic(NewRangeError("character index out of range"))
	}
	return ValueFromString(string(s.c[pos]))
}

// SetElem implements IArray.
func (s *String) SetElem(pos int, val Value) {
	panic(NewNotAssignableError(ValueFromString(s.s)))
}

// PushElem implements IArray.
func (s *String) PushElem(val Value) {
	panic(NewNotAssignableError(ValueFromString(s.s)))
}

func (s *String) initChars() {
	if s.c == nil {
		s.c = []rune(s.s)
	}
}

var _stringMethods map[string]Method

func init() {
	_stringMethods = map[string]Method{
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
