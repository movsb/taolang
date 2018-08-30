package main

import (
	"errors"
	"fmt"
)

func panicf(f string, args ...interface{}) {
	panic(fmt.Sprintf(f, args...))
}

func toErr(except interface{}) (err error) {
	switch typed := except.(type) {
	case nil:
		return nil
	case error:
		err = typed
	case string:
		err = errors.New(typed)
	default:
		err = fmt.Errorf("%v", typed)
	}
	return
}
