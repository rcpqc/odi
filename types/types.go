package types

import (
	"reflect"
	"sync"
)

const MaxKinds = 32

var (
	Any = reflect.TypeOf((*any)(nil)).Elem()

	Bool = reflect.TypeOf((*bool)(nil)).Elem()
	Byte = reflect.TypeOf((*byte)(nil)).Elem()

	Int   = reflect.TypeOf((*int)(nil)).Elem()
	Int8  = reflect.TypeOf((*int8)(nil)).Elem()
	Int16 = reflect.TypeOf((*int16)(nil)).Elem()
	Int32 = reflect.TypeOf((*int32)(nil)).Elem()
	Int64 = reflect.TypeOf((*int64)(nil)).Elem()

	Uint   = reflect.TypeOf((*uint)(nil)).Elem()
	Uint8  = reflect.TypeOf((*uint8)(nil)).Elem()
	Uint16 = reflect.TypeOf((*uint16)(nil)).Elem()
	Uint32 = reflect.TypeOf((*uint32)(nil)).Elem()
	Uint64 = reflect.TypeOf((*uint64)(nil)).Elem()

	Float32 = reflect.TypeOf((*float32)(nil)).Elem()
	Float64 = reflect.TypeOf((*float64)(nil)).Elem()

	String = reflect.TypeOf((*string)(nil)).Elem()
	Bytes  = reflect.TypeOf((*[]byte)(nil)).Elem()
)

var cache sync.Map

type typeKey struct {
	typ reflect.Type
	tag string
}

func LoadOrCreate(typ reflect.Type, tag string, constructor func() interface{}) (interface{}, bool) {
	tk := typeKey{typ, tag}
	if f, ok := cache.Load(tk); ok {
		return f.(func() interface{})(), true
	}
	var once sync.Once
	var res interface{}
	f, loaded := cache.LoadOrStore(tk, func() interface{} {
		once.Do(func() {
			res = constructor()
			cache.Store(tk, func() interface{} { return res })
		})
		return res
	})
	return f.(func() interface{})(), loaded
}
