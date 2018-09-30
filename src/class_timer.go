package main

import "time"

// Timer is a timer.
type Timer struct {
	timer *time.Timer
}

// NewTimer news a timer.
func NewTimer(ctx *Context, callback Value, timeout int) *Timer {
	t := time.NewTimer(time.Millisecond * time.Duration(timeout))
	go func() {
		select {
		case <-t.C:
			Async(func() {
				CallFunc(ctx, callback)
			})
		}
	}()
	return &Timer{timer: t}
}

func _NewTimer(this interface{}, ctx *Context, args *Values) Value {
	if args.Len() < 1 {
		panic(NewTypeError("Timer: callback expected"))
	} else if args.Len() < 2 {
		panic(NewTypeError("Timer: timeout expected"))
	}
	var callback = args.Shift()
	if !callback.isCallable() {
		panic(NewTypeError("Timer: callback must be a callable"))
	}
	timeout := args.Shift()
	if !timeout.isNumber() {
		panic(NewTypeError("Timer: timeout must be a number"))
	}
	t := NewTimer(ctx, callback, timeout.number())
	return ValueFromObject(t)
}

// Key implements KeyIndexer.
func (t *Timer) Key(key string) Value {
	if fn, ok := _timerMethods[key]; ok {
		return ValueFromBuiltin(t, key, fn)
	}
	return ValueFromNil()
}

// SetKey implements KeyIndexer.
func (t *Timer) SetKey(key string, val Value) {
	panic(NewNotAssignableError(ValueFromObject(t)))
}

var _timerMethods map[string]BuiltinFunction

func init() {
	_timerMethods = map[string]BuiltinFunction{
		"stop": _timerStop,
	}
}

func _timerStop(this interface{}, ctx *Context, args *Values) Value {
	t := this.(*Timer).timer
	return ValueFromBoolean(t.Stop())
}
