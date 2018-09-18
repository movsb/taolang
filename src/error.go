package main

import (
	"fmt"
)

// SyntaxError is
type SyntaxError struct {
	err string
}

// NewSyntaxError news
func NewSyntaxError(format string, args ...interface{}) SyntaxError {
	return SyntaxError{
		err: fmt.Sprintf(format, args...),
	}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s", e.err)
}

// NameError is
type NameError struct {
	err string
}

// NewNameError news
func NewNameError(format string, args ...interface{}) NameError {
	return NameError{
		err: fmt.Sprintf(format, args...),
	}
}

func (e NameError) Error() string {
	return fmt.Sprintf("NameError: %s", e.err)
}

// TypeError is
type TypeError struct {
	err string
}

// NewTypeError news
func NewTypeError(format string, args ...interface{}) TypeError {
	return TypeError{
		err: fmt.Sprintf(format, args...),
	}
}

func (e TypeError) Error() string {
	return fmt.Sprintf(
		"TypeError: %s", e.err,
	)
}

// NotCallableError is
type NotCallableError struct {
	value Value
}

// RangeError is
type RangeError struct {
	err string
}

// NewRangeError news
func NewRangeError(format string, args ...interface{}) RangeError {
	return RangeError{
		err: fmt.Sprintf(format, args...),
	}
}

func (e RangeError) Error() string {
	return fmt.Sprintf(
		"RangeError: %s", e.err,
	)
}

// NewNotCallableError news
func NewNotCallableError(value Value) NotCallableError {
	return NotCallableError{
		value: value,
	}
}

func (e NotCallableError) Error() string {
	return fmt.Sprintf("NotCallableError: %v (type: %s) is not callable", e.value, e.value.TypeName())
}

// NotIndexableError is
type NotIndexableError struct {
	value Value
}

// NewNotIndexableError news
func NewNotIndexableError(value Value) NotIndexableError {
	return NotIndexableError{
		value: value,
	}
}

func (e NotIndexableError) Error() string {
	return fmt.Sprintf("NotIndexableError: %v (type: %s) is not indexable", e.value, e.value.TypeName())
}

// NotAssignableError is
type NotAssignableError struct {
	value Value
}

// NewNotAssignableError news
func NewNotAssignableError(value Value) NotAssignableError {
	return NotAssignableError{
		value: value,
	}
}

func (e NotAssignableError) Error() string {
	return fmt.Sprintf("NotAssignableError: %v (type: %s) is not assignable", e.value, e.value.TypeName())
}

// KeyTypeError is
type KeyTypeError struct {
	value Value
}

// NewKeyTypeError news
func NewKeyTypeError(value Value) KeyTypeError {
	return KeyTypeError{
		value: value,
	}
}

func (e KeyTypeError) Error() string {
	return fmt.Sprintf(
		"KeyTypeError: cannot use %v (type: %s) as key",
		e.value, e.value.TypeName(),
	)
}
