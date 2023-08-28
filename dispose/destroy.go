package dispose

import (
	"context"
	"reflect"
	"strings"

	"github.com/rcpqc/odi/types"
)

var destoryers [types.MaxKinds]func(ctx context.Context, target reflect.Value) error

func init() {
	destoryers[reflect.Array] = destorySlice
	destoryers[reflect.Interface] = destoryInterface
	destoryers[reflect.Map] = destoryMap
	destoryers[reflect.Pointer] = destoryPointer
	destoryers[reflect.Slice] = destorySlice
	destoryers[reflect.Struct] = destoryStruct
}

func destoryPointer(ctx context.Context, target reflect.Value) error {
	if target.IsNil() {
		return nil
	}
	return destory(ctx, target.Elem())
}

func destorySlice(ctx context.Context, target reflect.Value) error {
	var anyerr error
	for i := 0; i < target.Len(); i++ {
		if err := destory(ctx, target.Index(i)); err != nil && anyerr == nil {
			anyerr = err
		}
	}
	return anyerr
}

func destoryMap(ctx context.Context, target reflect.Value) error {
	var anyerr error
	iter := target.MapRange()
	for iter.Next() {
		if err := destory(ctx, iter.Value()); err != nil && anyerr == nil {
			anyerr = err
		}
	}
	return anyerr
}

func destoryStruct(ctx context.Context, target reflect.Value) error {
	var anyerr error
	tagKey := ctxGetTagKey(ctx)
	for i := 0; i < target.NumField(); i++ {
		if !target.Field(i).CanSet() {
			continue
		}
		tag := target.Type().Field(i).Tag.Get(tagKey)
		tag = strings.Split(tag, ",")[0]
		if tag == "-" {
			continue
		}
		if err := destory(ctx, target.Field(i)); err != nil && anyerr == nil {
			anyerr = err
		}
	}
	return anyerr
}

func destoryInterface(ctx context.Context, target reflect.Value) error {
	if target.IsNil() {
		return nil
	}
	if iface, ok := target.Interface().(IDispose); ok {
		if err := iface.Dispose(); err != nil {
			return err
		}
	} else {
		if err := destory(ctx, target.Elem()); err != nil {
			return err
		}
	}
	return nil
}

func destory(ctx context.Context, target reflect.Value) error {
	destoryer := destoryers[target.Kind()]
	if destoryer == nil {
		return nil
	}
	return destoryer(ctx, target)
}
