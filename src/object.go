package main

import (
	"bytes"
	"fmt"
	"sort"
)

// KeyGetter is implemented by those who supports key getters.
type KeyGetter interface {
	GetKey(key string) Value
}

// KeySetter is implemented by those who supports key setters.
type KeySetter interface {
	SetKey(key string, val Value)
}

// ElemGetter is implemented by those who supports element getters.
type ElemGetter interface {
	Len() int
	GetElem(pos int) Value
}

// ElemSetter is implemented by those who supports element setters.
type ElemSetter interface {
	Len() int
	SetElem(pos int, val Value)
	PushElem(val Value)
}

// KeyAssigner is implemented by those who can be assigned.
type KeyAssigner interface {
	KeyAssign(key string, val Value)
}

// ElemAssigner is implemented by those who can be assigned.
type ElemAssigner interface {
	ElemAssign(elem int, val Value)
}

// Callable is a callable.
type Callable interface {
	Execute(ctx *Context, args *Values) Value
}

// Object is either an object or an array.
type Object struct {
	elems []Value          // array elements
	props map[string]Value // object properties
	array bool
}

// NewObject news an object.
func NewObject() *Object {
	o := &Object{}
	o.props = make(map[string]Value)
	return o
}

// NewArray news an array.
func NewArray(elems ...Value) *Object {
	o := NewObject()
	o.array = true
	o.elems = elems
	return o
}

// GetKey gets a value by key.
func (o *Object) GetKey(key string) Value {
	if o.array {
		if key == "length" {
			return ValueFromNumber(o.Len())
		}
	}
	if prop, ok := o.props[key]; ok {
		return prop
	}
	if fn, ok := _arrayMethods[key]; ok {
		return ValueFromBuiltin(o, key, fn)
	}
	return ValueFromNil()
}

// SetKey sets a value by key.
func (o *Object) SetKey(key string, val Value) {
	o.props[key] = val
}

// KeyAssign implements KeyAssigner.
func (o *Object) KeyAssign(key string, val Value) {
	o.SetKey(key, val)
}

// Len implements ElemGetter/ElemSetter.
func (o *Object) Len() int {
	return len(o.elems)
}

// GetElem implements ElemGetter.
func (o *Object) GetElem(pos int) Value {
	if pos < 0 || pos > len(o.elems)-1 {
		panic(NewRangeError("array index out of range"))
	}
	return o.elems[pos]
}

// SetElem implements ElemSetter.
func (o *Object) SetElem(pos int, val Value) {
	if pos < 0 || pos > len(o.elems)-1 {
		panic(NewRangeError("array index out of range"))
	}
	o.elems[pos] = val
}

// ElemAssign implements ElemAssigner.
func (o *Object) ElemAssign(elem int, val Value) {
	o.SetElem(elem, val)
}

// PushElem implements ElemSetter.
func (o *Object) PushElem(val Value) {
	o.elems = append(o.elems, val)
}

func (o *Object) sortedKeys() []string {
	keys := make([]string, len(o.props))
	i := 0
	for key := range o.props {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

func (o *Object) String() string {
	if o.array {
		buf := bytes.NewBuffer(nil)
		buf.WriteString("[")
		for i, n := 0, o.Len(); i < n; i++ {
			elem := o.GetElem(i)
			buf.WriteString(elem.String())
			if i != n-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("]")
		return buf.String()
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString("{")
	n := len(o.props)
	for i, key := range o.sortedKeys() {
		// TODO key may have invalid characters.
		buf.WriteString(fmt.Sprintf(`%s:%v`, key, o.props[key]))
		if i != n-1 {
			buf.WriteString(",")
		}
		i++
	}
	buf.WriteString("}")
	return buf.String()
}

/// Array function implementations below
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array
// Javascript array methods.

var _arrayMethods map[string]BuiltinFunction

func init() {
	_arrayMethods = map[string]BuiltinFunction{
		"each":    _arrayEach,
		"filter":  _arrayFilter,
		"find":    _arrayFind,
		"groupBy": _arrayGroupBy,
		"join":    _arrayJoin,
		"map":     _arrayMap,
		"push":    _arrayPush,
		"pop":     _arrayPop,
		"reduce":  _arrayReduce,
		"select":  _arraySelect,
		"splice":  _arraySplice,
		"unshift": _arrayUnshift,
		"where":   _arrayWhere,
	}
}

// Splice changes the contents of an array by removing existing elements and/or adding new elements.
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/splice
func _arraySplice(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	start := 0
	if args.Len() < 1 || !args.At(0).isNumber() {
		panic(NewTypeError("splice: start must be number"))
	}
	start = args.Shift().number()
	if start > o.Len() {
		start = o.Len()
	} else if start < 0 {
		if -start > o.Len() {
			start = 0
		} else {
			start += o.Len()
		}
	}
	deleteCount := 0
	if args.Len() >= 1 {
		if !args.At(0).isNumber() {
			panic(NewTypeError("splice: deleteCount must be number"))
		}
		deleteCount = args.Shift().number()
		if deleteCount > o.Len()-start {
			deleteCount = o.Len() - start
		}
		if deleteCount <= 0 {

		}
	} else {
		deleteCount = o.Len() - start
	}
	deletedElements := []Value{}
	if deleteCount > 0 {
		deletedElements = make([]Value, deleteCount)
		copy(deletedElements, o.elems[start:start+deleteCount])
		o.elems = append(o.elems[0:start], o.elems[start+deleteCount:]...)
	}
	if args.Len() > 0 {
		elems := make([]Value, len(o.elems)+args.Len())
		copy(elems, o.elems[:start])
		copy(elems[start:], args.values)
		copy(elems[start+args.Len():], o.elems[start:])
		o.elems = elems
	}
	return ValueFromObject(NewArray(deletedElements...))
}

// Unshift adds elements to the beginning of the array and returns the new length of the array.
// https://github.com/golang/go/wiki/SliceTricks#push-frontunshift
func _arrayUnshift(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	for _, v := range args.values {
		o.elems = append([]Value{v}, o.elems...)
	}
	return ValueFromNumber(o.Len())
}

// Push adds one or more elements to the end of an array and returns the new length of the array.
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/push
func _arrayPush(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	o.elems = append(o.elems, args.values...)
	return ValueFromNumber(o.Len())
}

// Pop removes the last element from an array and returns that element.
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/pop
func _arrayPop(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	if o.Len() > 0 {
		value := o.elems[o.Len()-1]
		o.elems = o.elems[:o.Len()-1]
		return value
	}
	return ValueFromNil()
}

// Join joins all elements of an array into a string and returns this string.
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/join
func _arrayJoin(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	if o.Len() <= 0 {
		return ValueFromString("")
	}
	sep := ""
	if args.Len() >= 1 {
		sep = fmt.Sprint(args.At(0))
	}
	n := o.Len()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < n-1; i++ {
		s := fmt.Sprintf("%v%s", o.elems[i], sep)
		buf.WriteString(s)
	}
	buf.WriteString(fmt.Sprint(o.elems[n-1]))
	return ValueFromString(buf.String())
}

/// functional methods implementations below.

func _arrayCall(ctx *Context, lambda Value, args ...Value) Value {
	ctx = NewContext("--lambda--", nil)
	if !lambda.isCallable() {
		panic(NewNotCallableError(lambda))
	}
	return lambda.callable().Execute(ctx, NewValues(args...))
}

// Each iterates each element of the array and invokes callback.
func (o *Object) Each(callback func(elem Value, index Value) bool) {
	for i, n, next := 0, o.Len(), true; i < n && next; i++ {
		next = callback(o.elems[i], ValueFromNumber(i))
	}
}

// Each iterates over a list of elements, yielding each in turn to an iteratee function.
func _arrayEach(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	object := ValueFromObject(o)
	o.Each(func(elem Value, index Value) bool {
		_arrayCall(ctx, args.At(0), elem, index, object)
		return true
	})
	return ValueFromNil()
}

// Map produces a new array of values by mapping each value.
func _arrayMap(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	object := ValueFromObject(o)
	values := make([]Value, 0, o.Len())
	o.Each(func(elem Value, index Value) bool {
		data := _arrayCall(ctx, args.At(0), elem, index, object)
		values = append(values, data)
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// Reduce boils down the array into a single value.
func _arrayReduce(this interface{}, ctx *Context, args *Values) Value {
	if args.Len() < 2 {
		// TODO arguments error
		panic(NewTypeError("usage: reduce(lambda, init)"))
	}
	o := this.(*Object)
	object := ValueFromObject(o)
	memo := args.At(1)
	o.Each(func(elem Value, index Value) bool {
		memo = _arrayCall(ctx, args.At(0), memo, elem, index, object)
		return true
	})
	return memo
}

// Find finds the first value.
func _arrayFind(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	found := Value{}
	o.Each(func(elem Value, index Value) bool {
		if _arrayCall(ctx, args.At(0), elem).Truth(ctx) {
			found = elem
			return false
		}
		return true
	})
	return found
}

// Filter filters values.
func _arrayFilter(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	values := make([]Value, 0, o.Len())
	o.Each(func(elem Value, index Value) bool {
		if _arrayCall(ctx, args.At(0), elem).Truth(ctx) {
			values = append(values, elem)
		}
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// Where filters objects by column conditions.
// same as Filter currently.
func _arrayWhere(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	values := make([]Value, 0, o.Len())
	o.Each(func(elem Value, index Value) bool {
		if _arrayCall(ctx, args.At(0), elem).Truth(ctx) {
			values = append(values, elem)
		}
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// Select selects fields as array.
func _arraySelect(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	values := make([]Value, 0, o.Len())
	o.Each(func(elem Value, index Value) bool {
		value := _arrayCall(ctx, args.At(0), elem)
		values = append(values, value)
		return true
	})
	return ValueFromObject(NewArray(values...))
}

// GroupBy groups objects by property.
func _arrayGroupBy(this interface{}, ctx *Context, args *Values) Value {
	o := this.(*Object)
	maps := make(map[Value][]Value)
	keys := make([]Value, 0) // make map output ordered
	o.Each(func(elem Value, index Value) bool {
		key := _arrayCall(ctx, args.At(0), elem)
		if _, ok := maps[key]; !ok {
			keys = append(keys, key)
		}
		maps[key] = append(maps[key], elem)
		return true
	})
	group := NewArray()
	for _, key := range keys {
		obj := NewArray()
		obj.SetKey("group", key)
		obj.elems = maps[key]
		group.PushElem(ValueFromObject(obj))
	}
	return ValueFromObject(group)
}
