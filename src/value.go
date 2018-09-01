package main

import (
	"bytes"
	"fmt"
)

// ValueType is the type of a value.
type ValueType int

const (
	vtNil ValueType = iota
	vtBoolean
	vtNumber
	vtString
	vtVariable
	vtObject
	vtFunction
	vtBuiltin
)

var typeNames = map[ValueType]string{
	vtNil:      "nil",
	vtBoolean:  "boolean",
	vtNumber:   "number",
	vtString:   "string",
	vtVariable: "variable",
	vtObject:   "object",
	vtFunction: "function",
	vtBuiltin:  "builtin",
}

// Value holds a union(dynamic) value identified by Type.
type Value struct {
	Type  ValueType
	value interface{}
}

// ValueFromNil creates a nil value.
func ValueFromNil() Value {
	return Value{
		Type: vtNil,
	}
}

// ValueFromBoolean creates a boolean value.
func ValueFromBoolean(b bool) Value {
	return Value{
		Type:  vtBoolean,
		value: b,
	}
}

// ValueFromNumber creates a number value.
func ValueFromNumber(num int) Value {
	return Value{
		Type:  vtNumber,
		value: num,
	}
}

// ValueFromString creates a string value.
func ValueFromString(str string) Value {
	return Value{
		Type:  vtString,
		value: str,
	}
}

// ValueFromVariable creates a variable value.
// It references a variable by its name.
func ValueFromVariable(name string) Value {
	return Value{
		Type:  vtVariable,
		value: name,
	}
}

// ValueFromObject creates an array | object value.
func ValueFromObject(obj *Object) Value {
	return Value{
		Type:  vtObject,
		value: obj,
	}
}

// ValueFromFunction creates a evaluated function expression value.
func ValueFromFunction(fn *EvaluatedFunctionExpression) Value {
	return Value{
		Type:  vtFunction,
		value: fn,
	}
}

// ValueFromBuiltin creates a builtin function value.
func ValueFromBuiltin(name string, fn func(*Context, *Values) Value) Value {
	return Value{
		Type: vtBuiltin,
		value: &Builtin{
			name: name,
			fn:   fn,
		},
	}
}

func (v *Value) isNil() bool {
	return v.Type == vtNil
}

func (v *Value) isBoolean() bool {
	return v.Type == vtBoolean
}

func (v *Value) isNumber() bool {
	return v.Type == vtNumber
}

func (v *Value) isString() bool {
	return v.Type == vtString
}

func (v *Value) isObject() bool {
	return v.Type == vtObject
}

func (v *Value) isVariable() bool {
	return v.Type == vtVariable
}

func (v *Value) isFunction() bool {
	return v.Type == vtFunction
}

func (v *Value) isBuiltin() bool {
	return v.Type == vtBuiltin
}

func (v *Value) checkType(vt ValueType) {
	if v.Type != vt {
		panic("wrong use")
	}
}

func (v *Value) boolean() bool {
	v.checkType(vtBoolean)
	return v.value.(bool)
}

func (v *Value) number() int {
	v.checkType(vtNumber)
	return v.value.(int)
}

func (v *Value) str() string {
	v.checkType(vtString)
	return v.value.(string)
}

func (v *Value) variable() string {
	v.checkType(vtVariable)
	return v.value.(string)
}

func (v *Value) object() *Object {
	v.checkType(vtObject)
	return v.value.(*Object)
}

func (v *Value) function() *EvaluatedFunctionExpression {
	v.checkType(vtFunction)
	return v.value.(*EvaluatedFunctionExpression)
}

func (v *Value) builtin() *Builtin {
	v.checkType(vtBuiltin)
	return v.value.(*Builtin)
}

// Evaluate implements Expression.
func (v Value) Evaluate(ctx *Context) Value {
	switch v.Type {
	case vtNil, vtBoolean, vtNumber, vtString:
		return v
	case vtVariable:
		return ctx.MustFind(v.variable(), true)
	case vtObject:
		return v
	case vtFunction:
		return v
	case vtBuiltin:
		return v
	default:
		panic("cannot evaluate value on type")
	}
}

// Address implements Addresser.
func (v Value) Address(ctx *Context) *Value {
	if v.isVariable() {
		return ctx.GetAddress(v.variable())
	}
	panicf("not assignable: %v (type: %s)", v.value, v.TypeName())
	return nil
}

// TypeName returns the value type as string.
func (v Value) TypeName() string {
	return typeNames[v.Type]
}

func (v Value) String() string {
	switch v.Type {
	case vtNil:
		return "nil"
	case vtBoolean:
		return fmt.Sprint(v.boolean())
	case vtNumber:
		return fmt.Sprint(v.number())
	case vtString:
		return `"` + v.str() + `"`
	case vtFunction:
		expr := v.function()
		name := expr.expr.name
		if name == "" {
			name = "<anonymous>"
		}
		return fmt.Sprintf("function(%s)", name)
	case vtBuiltin:
		return fmt.Sprintf("builtin(%s)", v.builtin().name)
	case vtObject:
		if !v.object().IsArray() {
			buf := bytes.NewBuffer(nil)
			buf.WriteString("{")
			n := len(v.object().props)
			i := 0
			for k, p := range v.object().props {
				// TODO k may have invalid characters.
				buf.WriteString(fmt.Sprintf(`%s:%v`, k, p))
				if i != n-1 {
					buf.WriteString(",")
				}
				i++
			}
			buf.WriteString("}")
			return buf.String()
		} else {
			buf := bytes.NewBuffer(nil)
			buf.WriteString("[")
			for i, n := 0, v.object().Len(); i < n; i++ {
				elem := v.object().Elem(i)
				buf.WriteString(elem.String())
				if i != n-1 {
					buf.WriteString(",")
				}
			}
			buf.WriteString("]")
			return buf.String()
		}
	}
	return fmt.Sprintf("unknown value")
}

// Truth return true if value represents a true value.
// A value is considered true when:
func (v Value) Truth(ctx *Context) bool {
	switch v.Type {
	case vtNil:
		return false
	case vtNumber:
		return v.number() != 0
	case vtString:
		return v.str() != ""
	case vtBoolean:
		return v.boolean()
	case vtFunction, vtBuiltin:
		return true
	case vtVariable:
		return ctx.MustFind(v.variable(), true).Truth(ctx)
	case vtObject:
		obj := v.object()
		if obj.IsArray() {
			return obj.Len() > 0
		}
		return len(obj.props) > 0
	}
	panicf("unknown truth type")
	return false
}

// Values is a collection of values.
type Values struct {
	values []Value
}

// NewValues news a Values.
func NewValues(values ...Value) *Values {
	v := &Values{}
	for _, value := range values {
		v.values = append(v.values, value)
	}
	return v
}

// At returns
func (v *Values) At(i int) Value {
	if i < 0 && i > v.Len()-1 {
		panic("Values' index out of range")
	}
	return v.values[i]
}

// Len lens the values.
func (v *Values) Len() int {
	return len(v.values)
}

// All alls the values.
func (v *Values) All() []interface{} {
	var i []interface{}
	for _, value := range v.values {
		i = append(i, value)
	}
	return i
}
