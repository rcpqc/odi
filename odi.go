package odi

import (
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/container"
	"github.com/rcpqc/odi/resolve"
)

// Provide register a kind of object
func Provide(kind string, constructor func() any) {
	container.Provide(kind, constructor)
}

// Resolve parse and construct an new object by source data
func Resolve(source any, opts ...resolve.Option) (any, error) {
	return resolve.Invoke(source, opts...)
}

// Iterate traverse the object and callback the fields that implement the specified interface
func Iterate[T any](object any, cb func(iface T) error) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Interface {
		return fmt.Errorf("must be a generic implementation of the interface")
	}
	return iterate(reflect.ValueOf(object), reflect.TypeOf((*T)(nil)).Elem(),
		func(target reflect.Value) error {
			return cb(target.Interface().(T))
		},
	)
}
