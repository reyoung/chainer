package chainer

import (
	"errors"
	"fmt"
	"reflect"
)

type ResultHolder struct {
	value interface{}
	error error
}

func Wrap(instance interface{}) *ResultHolder {
	return &ResultHolder{
		value: instance,
	}
}

func Call(method interface{}, args ...interface{}) *ResultHolder {
	rMethod := reflect.TypeOf(method)
	if rMethod.Kind() != reflect.Func {
		return &ResultHolder{
			error: errors.New("method should be function"),
		}
	}

	if rMethod.NumIn() != len(args) {
		return &ResultHolder{
			error: errors.New("num of args mismatch"),
		}
	}

	rArgs := make([]reflect.Value, rMethod.NumIn())

	for offset, arg := range args {
		if rh, ok := arg.(*ResultHolder); ok {
			if rh.error != nil {
				return &ResultHolder{error: rh.error}
			}
			arg = rh.value
		}
		aVal := reflect.ValueOf(arg)

		if !aVal.Type().ConvertibleTo(rMethod.In(offset)) {
			return &ResultHolder{error: fmt.Errorf("arg %d type mismatch", offset)}
		}
		rArgs[offset] = aVal.Convert(rMethod.In(offset))
	}

	results := reflect.ValueOf(method).Call(rArgs)
	switch len(results) {
	case 0:
		return &ResultHolder{}
	case 1:
		// is error
		err, ok := results[0].Interface().(error)
		if ok {
			return &ResultHolder{error: err}
		} else {
			return &ResultHolder{value: results[0].Interface()}
		}

	case 2:
		return &ResultHolder{value: results[0].Interface(), error: results[1].Interface().(error)}
	default:
		panic("not support")
	}
}

func (r *ResultHolder) Error() error {
	return r.error
}

func (r *ResultHolder) Value() interface{} {
	return r.value
}

func (r *ResultHolder) MustValue() interface{} {
	if r.error != nil {
		panic(r.error)
	}
	return r.value
}

func (r *ResultHolder) Then(callback interface{}) *ResultHolder {
	return Call(callback, r)
}
