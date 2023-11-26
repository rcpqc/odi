package odi

import (
	"reflect"
)

// IDispose custom dispose interface for an object
type IDispose interface{ Dispose() }

var IDisposeType = reflect.TypeOf((*IDispose)(nil)).Elem()

func dispose(target reflect.Value) {
	if !target.IsValid() {
		return
	}
	if target.Type().Implements(IDisposeType) {
		target.Interface().(IDispose).Dispose()
	}
	switch target.Kind() {
	case reflect.Pointer:
		if !target.IsNil() {
			dispose(target.Elem())
		}
	case reflect.Struct:
		for i := 0; i < target.NumField(); i++ {
			if target.Field(i).CanSet() {
				dispose(target.Field(i))
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < target.Len(); i++ {
			dispose(target.Index(i))
		}
	case reflect.Map:
		for iter := target.MapRange(); iter.Next(); {
			dispose(iter.Key())
			dispose(iter.Value())
		}
	case reflect.Interface:
		dispose(target.Elem())
	}
}
