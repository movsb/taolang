package main

import (
	"bytes"
	"fmt"
)

type ValueType int

const (
	vtNil ValueType = iota
	vtNumber
	vtString
	vtBoolean
	vtFunction
	vtVariable
	vtBuiltin
	vtObject
)

// how to make this a c-like union struct ?
type Value struct {
	Type     ValueType
	Bool     bool
	Number   int
	Str      string
	Func     Expression
	Variable string
	Builtin  Builtin
	Object   *Object
}

func ValueFromNil() *Value {
	v := Value{}
	v.SetNil()
	return &v
}

func ValueFromBoolean(b bool) *Value {
	v := Value{}
	v.SetBoolean(b)
	return &v
}

func ValueFromNumber(num int) *Value {
	v := Value{}
	v.SetNumber(num)
	return &v
}

func ValueFromString(str string) *Value {
	v := Value{}
	v.SetString(str)
	return &v
}

func ValueFromFunction(name string, expr Expression) *Value {
	v := Value{}
	v.Str = name
	v.SetFunction(expr)
	return &v
}

func ValueFromVariable(name string) *Value {
	v := Value{}
	v.SetVariable(name)
	return &v
}

func ValueFromBuiltin(name string, builtin Builtin) *Value {
	v := Value{}
	v.Str = name
	v.SetBuiltin(builtin)
	return &v
}

func ValueFromObject(obj *Object) *Value {
	v := Value{}
	v.SetObject(obj)
	return &v
}

func (v *Value) SetNil() {
	v.Type = vtNil
}

func (v *Value) SetBoolean(b bool) {
	v.Type = vtBoolean
	v.Bool = b
}

func (v *Value) SetNumber(num int) {
	v.Type = vtNumber
	v.Number = num
}

func (v *Value) SetString(str string) {
	v.Type = vtString
	v.Str = str
}

func (v *Value) SetFunction(expr Expression) {
	v.Type = vtFunction
	v.Func = expr
}

func (v *Value) SetVariable(name string) {
	v.Type = vtVariable
	v.Variable = name
}

func (v *Value) SetBuiltin(builtin Builtin) {
	v.Type = vtBuiltin
	v.Builtin = builtin
}

func (v *Value) SetObject(obj *Object) {
	v.Type = vtObject
	v.Object = obj
}

func (v *Value) Evaluate(ctx *Context) *Value {
	switch v.Type {
	case vtNil, vtBoolean, vtNumber, vtString:
		cp := *v
		return &cp
	case vtFunction:
		return v
	case vtVariable:
		value := ctx.FindValue(v.Variable, true)
		if value == nil {
			panic(fmt.Sprintf("undefined symbol: %s", v.Variable))
		}
		return value
	case vtBuiltin:
		return v
	case vtObject:
		return v
	default:
		panic("cannot evaluate value on type")
	}
}

func (v *Value) String() string {
	switch v.Type {
	case vtNil:
		return "nil"
	case vtBoolean:
		return fmt.Sprint(v.Bool)
	case vtNumber:
		return fmt.Sprint(v.Number)
	case vtString:
		return v.Str
	case vtFunction:
		expr := v.Func.(*FunctionExpression)
		name := expr.name
		if name == "" {
			name = "<anonymous>"
		}
		return fmt.Sprintf("function(%s)", name)
	case vtBuiltin:
		return fmt.Sprintf("builtin(%s)", v.Str)
	case vtObject:
		if !v.Object.IsArray() {
			return fmt.Sprintf(`"[object]"`)
		} else {
			buf := bytes.NewBuffer(nil)
			buf.WriteString("[")
			for i, n := 0, v.Object.Len(); i < n; i++ {
				buf.WriteString(v.Object.Elem(i).String())
				if i != n-1 {
					buf.WriteString(",")
				}
			}
			buf.WriteString("]")
			return buf.String()
		}
	}
	return fmt.Sprintf("unknown(%p)", v)
}

func (v *Value) Interface() interface{} {
	switch v.Type {
	case vtNil:
		return nil
	case vtNumber:
		return v.Number
	case vtString:
		return v.Str
	case vtBoolean:
		return v.Bool
	case vtFunction:
		return v.Func
	case vtVariable:
		return v.Variable
	case vtBuiltin:
		return v.Builtin
	case vtObject:
		return v.Object
	default:
		return nil
	}
}

func (v *Value) Truthy(ctx *Context) bool {
	switch v.Type {
	case vtNil:
		return false
	case vtNumber:
		return v.Number != 0
	case vtString:
		return v.Str != ""
	case vtBoolean:
		return v.Bool
	case vtFunction, vtBuiltin, vtObject:
		return true
	case vtVariable:
		value := ctx.FindValue(v.Variable, true)
		if value == nil {
			panicf("variable is not defined: %s", v.Variable)
		}
		return value.Truthy(ctx)
	}
	panicf("unknown truthy type")
	return false
}

type Values []*Value

func (v *Values) Len() int {
	return len(*v)
}

func (v *Values) ToInterfaces() []interface{} {
	var i []interface{}
	for _, value := range *v {
		i = append(i, value)
	}
	return i
}
