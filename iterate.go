package odi

import (
	"reflect"

	"github.com/rcpqc/odi/types"
)

var iterators [types.MaxKinds]func(target reflect.Value, iface reflect.Type, cb func(reflect.Value) error) error

func init() {
	iterators[reflect.Array] = iterateSlice
	iterators[reflect.Interface] = iterateInterface
	iterators[reflect.Map] = iterateMap
	iterators[reflect.Pointer] = iteratePointer
	iterators[reflect.Slice] = iterateSlice
	iterators[reflect.Struct] = iterateStruct
}

func iteratePointer(target reflect.Value, iface reflect.Type, cb func(reflect.Value) error) error {
	if target.IsNil() {
		return nil
	}
	return iterate(target.Elem(), iface, cb)
}

func iterateSlice(target reflect.Value, iface reflect.Type, cb func(reflect.Value) error) error {
	for i := 0; i < target.Len(); i++ {
		if err := iterate(target.Index(i), iface, cb); err != nil {
			return err
		}
	}
	return nil
}

func iterateMap(target reflect.Value, iface reflect.Type, cb func(reflect.Value) error) error {
	iter := target.MapRange()
	for iter.Next() {
		if err := iterate(iter.Key(), iface, cb); err != nil {
			return err
		}
		if err := iterate(iter.Value(), iface, cb); err != nil {
			return err
		}
	}
	return nil
}

func iterateStruct(target reflect.Value, iface reflect.Type, cb func(reflect.Value) error) error {
	for i := 0; i < target.NumField(); i++ {
		if !target.Field(i).CanSet() {
			continue
		}
		if err := iterate(target.Field(i), iface, cb); err != nil {
			return err
		}
	}
	return nil
}

func iterateInterface(target reflect.Value, iface reflect.Type, cb func(reflect.Value) error) error {
	return iterate(target.Elem(), iface, cb)
}

func iterate(target reflect.Value, iface reflect.Type, cb func(reflect.Value) error) error {
	if !target.IsValid() {
		return nil
	}
	if target.Type().Implements(iface) {
		return cb(target)
	}
	if iterator := iterators[target.Kind()]; iterator != nil {
		return iterator(target, iface, cb)
	}
	return nil
}
