package resolve

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/container"
	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

func invoke(ctx context.Context, src reflect.Value) (any, error) {
	kind, err := classify(ctx, src)
	if err != nil {
		return nil, errs.Newf("classify -> %v", err)
	}
	object, err := container.Make(kind)
	if err != nil {
		return nil, errs.Newf("container.make -> %v", err)
	}
	if err := inject(ctx, reflect.ValueOf(object), src); err != nil {
		return nil, err
	}
	return object, nil
}

func classify(ctx context.Context, src reflect.Value) (string, error) {
	if src.Kind() != reflect.Map {
		return "", fmt.Errorf("expect map but %v", src.Kind())
	}
	key := ctxGetObjectKey(ctx)
	rkind := src.MapIndex(reflect.ValueOf(key))
	if !rkind.IsValid() {
		return "", fmt.Errorf("not found object_key(%s)", key)
	}
	if rkind.Type() == types.Any {
		rkind = rkind.Elem()
	}
	if rkind.Type() != types.String {
		return "", fmt.Errorf("illegal object_kind(%v)", rkind.Interface())
	}
	return rkind.Interface().(string), nil
}
