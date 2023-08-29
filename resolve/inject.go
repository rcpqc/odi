package resolve

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/errs"

	"github.com/rcpqc/odi/types"
)

// IResolve custom resolve interface for an object
type IResolve interface{ Resolve(src any) error }

var IResolveType = reflect.TypeOf((*IResolve)(nil)).Elem()

var injectors [types.MaxKinds]func(ctx context.Context, dst, src reflect.Value) error

func init() {
	injectors[reflect.Bool] = injectBool
	injectors[reflect.Int] = injectInt
	injectors[reflect.Int8] = injectInt
	injectors[reflect.Int16] = injectInt
	injectors[reflect.Int32] = injectInt
	injectors[reflect.Int64] = injectInt
	injectors[reflect.Uint] = injectUint
	injectors[reflect.Uint8] = injectUint
	injectors[reflect.Uint16] = injectUint
	injectors[reflect.Uint32] = injectUint
	injectors[reflect.Uint64] = injectUint
	injectors[reflect.Uintptr] = injectUint
	injectors[reflect.Float32] = injectFloat
	injectors[reflect.Float64] = injectFloat
	injectors[reflect.Array] = injectArray
	injectors[reflect.Interface] = injectInterface
	injectors[reflect.Map] = injectMap
	injectors[reflect.Pointer] = injectPointer
	injectors[reflect.Slice] = injectSlice
	injectors[reflect.String] = injectString
	injectors[reflect.Struct] = injectStruct
}

func inject(ctx context.Context, dst, src reflect.Value) error {
	if src.IsValid() && src.Type() == types.Any {
		src = src.Elem()
	}
	if !src.IsValid() {
		return nil
	}
	injector := injectors[dst.Kind()]
	if injector == nil {
		return errs.Newf("not support kind: %v", dst.Kind())
	}
	err := injector(ctx, dst, src)
	if dst.CanAddr() && reflect.PointerTo(dst.Type()).Implements(IResolveType) {
		err = dst.Addr().Interface().(IResolve).Resolve(src.Interface())
	}
	return err
}

func injectPointer(ctx context.Context, dst, src reflect.Value) error {
	if dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}
	return inject(ctx, dst.Elem(), src)
}

func injectArray(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Array && src.Kind() != reflect.Slice {
		return errs.Newf("expect slice or array but %v", src.Kind())
	}
	dst.Set(reflect.New(dst.Type()).Elem())
	for i := 0; i < src.Len() && i < dst.Len(); i++ {
		elem := reflect.New(dst.Type().Elem()).Elem()
		if err := inject(ctx, elem, src.Index(i)); err != nil {
			return errs.New(err).Prefix(fmt.Sprintf("[%d]", i))
		}
		dst.Index(i).Set(elem)
	}
	return nil
}

func injectSlice(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Array && src.Kind() != reflect.Slice {
		return errs.Newf("expect slice or array but %v", src.Kind())
	}
	dst.Set(reflect.MakeSlice(dst.Type(), src.Len(), src.Cap()))
	for i := 0; i < src.Len(); i++ {
		elem := reflect.New(dst.Type().Elem()).Elem()
		if err := inject(ctx, elem, src.Index(i)); err != nil {
			return errs.New(err).Prefix(fmt.Sprintf("[%d]", i))
		}
		dst.Index(i).Set(elem)
	}
	return nil
}

func injectMap(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Map {
		return errs.Newf("expect map but %v", src.Kind())
	}
	dst.Set(reflect.MakeMap(dst.Type()))
	iter := src.MapRange()
	for iter.Next() {
		srcKey, srcVal := iter.Key(), iter.Value()
		dstKey, dstVal := reflect.New(dst.Type().Key()).Elem(), reflect.New(dst.Type().Elem()).Elem()
		if err := inject(ctx, dstKey, srcKey); err != nil {
			return errs.New(err).Prefix(fmt.Sprintf("[%s]", dstKey.String()))
		}
		if err := inject(ctx, dstVal, srcVal); err != nil {
			return errs.New(err).Prefix(fmt.Sprintf("[%s]", dstKey.String()))
		}
		dst.SetMapIndex(dstKey, dstVal)
	}
	return nil
}

func injectInterface(ctx context.Context, dst, src reflect.Value) error {
	if dst.Type() == types.Any && dst.CanSet() {
		dst.Set(src)
		return nil
	}
	object, err := invoke(ctx, src)
	if err != nil {
		return err
	}
	if !reflect.ValueOf(object).Type().Implements(dst.Type()) {
		return errs.Newf("the injected object does not implement the interface(%v)", dst.Type())
	}
	dst.Set(reflect.ValueOf(object))
	return nil
}
