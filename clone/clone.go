package clone

import (
	"reflect"

	"github.com/rcpqc/odi/types"
)

var Clone = clone

func cloneString(rv reflect.Value) reflect.Value {
	return reflect.ValueOf(rv.String())
}

func clonePtr(rv reflect.Value) reflect.Value {
	if rv.IsNil() {
		return rv
	}
	ptr := reflect.New(rv.Type().Elem())
	ptr.Elem().Set(clone(rv.Elem()))
	return ptr
}

func cloneArray(rv reflect.Value) reflect.Value {
	arr := reflect.New(rv.Type()).Elem()
	for i := 0; i < rv.Len(); i++ {
		arr.Index(i).Set(clone(rv.Index(i)))
	}
	return arr
}

func cloneSlice(rv reflect.Value) reflect.Value {
	if rv.IsNil() {
		return rv
	}
	slc := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Cap())
	for i := 0; i < rv.Len(); i++ {
		slc.Index(i).Set(clone(rv.Index(i)))
	}
	return slc
}

func cloneMap(rv reflect.Value) reflect.Value {
	if rv.IsNil() {
		return rv
	}
	m := reflect.MakeMapWithSize(rv.Type(), rv.Len())
	iter := rv.MapRange()
	for iter.Next() {
		m.SetMapIndex(clone(iter.Key()), clone(iter.Value()))
	}
	return m
}

func cloneStruct(rv reflect.Value) reflect.Value {
	stu := reflect.New(rv.Type()).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).CanSet() {
			stu.Field(i).Set(clone(rv.Field(i)))
		}
	}
	return stu
}

func cloneInterface(rv reflect.Value) reflect.Value {
	if rv.IsZero() {
		return rv
	}
	data := reflect.New(rv.Elem().Type()).Elem()
	data.Set(clone(rv.Elem()))
	return data
}

func clone(rv reflect.Value) reflect.Value {
	if rv.IsValid() && rv.CanInterface() {
		if iface, ok := rv.Interface().(types.Clonable); ok {
			return reflect.ValueOf(iface.OnClone())
		}
	}
	switch rv.Kind() {
	case reflect.String:
		return cloneString(rv)
	case reflect.Ptr:
		return clonePtr(rv)
	case reflect.Slice:
		return cloneSlice(rv)
	case reflect.Map:
		return cloneMap(rv)
	case reflect.Struct:
		return cloneStruct(rv)
	case reflect.Interface:
		return cloneInterface(rv)
	case reflect.Array:
		return cloneArray(rv)
	default:
		return rv
	}
}
