package dispose

import (
	"reflect"
	"strings"
)

// IDispose custom dispose interface for an object
type IDispose interface{ Dispose() error }

var Dispose = dispose

func disposePtr(rv reflect.Value) error {
	if rv.IsNil() {
		return nil
	}
	return dispose(rv.Elem())
}

func disposeSlice(rv reflect.Value) error {
	var anyerr error
	for i := 0; i < rv.Len(); i++ {
		if err := dispose(rv.Index(i)); err != nil && anyerr == nil {
			anyerr = err
		}
	}
	return anyerr
}

func disposeMap(rv reflect.Value) error {
	var anyerr error
	iter := rv.MapRange()
	for iter.Next() {
		if err := dispose(iter.Value()); err != nil && anyerr == nil {
			anyerr = err
		}
	}
	return anyerr
}

func disposeStruct(rv reflect.Value) error {
	var anyerr error
	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).CanSet() {
			continue
		}
		tag := rv.Type().Field(i).Tag.Get("odi")
		tag = strings.Split(tag, ",")[0]
		if tag == "-" {
			continue
		}
		if err := dispose(rv.Field(i)); err != nil && anyerr == nil {
			anyerr = err
		}
	}
	return anyerr
}

func disposeInterface(rv reflect.Value) error {
	var anyerr error
	if rv.IsNil() {
		return anyerr
	}
	if err := dispose(rv.Elem()); err != nil {
		anyerr = err
	}
	iface, ok := rv.Interface().(IDispose)
	if !ok {
		return anyerr
	}
	if err := iface.Dispose(); err != nil && anyerr == nil {
		anyerr = err
	}
	return anyerr
}

func dispose(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Ptr:
		return disposePtr(rv)
	case reflect.Array, reflect.Slice:
		return disposeSlice(rv)
	case reflect.Map:
		return disposeMap(rv)
	case reflect.Struct:
		return disposeStruct(rv)
	case reflect.Interface:
		return disposeInterface(rv)
	default:
		return nil
	}
}
