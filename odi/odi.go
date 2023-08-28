package odi

import (
	"reflect"

	"github.com/rcpqc/odi/container"
	"github.com/rcpqc/odi/dispose"
	"github.com/rcpqc/odi/resolve"
)

// Provide register a kind of object
func Provide(kind string, constructor func() any) {
	container.Provide(kind, constructor)
}

// Resolve parse and construct an new object by source data
func Resolve(data any, opts ...resolve.Option) (any, error) {
	opts = append(opts, resolve.WithKey("object"))
	return resolve.Invoke(data, opts...)
}

// Dispose release an object and callback IDispose interface
func Dispose(object any) error {
	return dispose.Dispose(reflect.ValueOf(object))
}

type IResolve = resolve.IResolve

type IDispose = dispose.IDispose
