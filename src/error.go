package main

import (
	"fmt"
)

// NotCallableError is
type NotCallableError struct {
	value Value
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
