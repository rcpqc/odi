package odi

import (
	"reflect"

	"github.com/rcpqc/odi/container"
	"github.com/rcpqc/odi/resolve"
)

// Provide register a kind of object
func Provide(kind string, constructor func() any) {
	container.Bind(kind, constructor)
}

// Resolve parse and construct an new object by source data
func Resolve(source any, opts ...resolve.Option) (any, error) {
	return resolve.Invoke(source, opts...)
}

// Dispose an object
func Dispose(target any) {
	dispose(reflect.ValueOf(target))
}
