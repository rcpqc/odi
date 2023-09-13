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
		return nil, errs.Newf("classify err: %v", err)
	}
	object, err := container.Create(kind)
	if err != nil {
		return nil, errs.Newf("container create err: %v", err)
	}
	if err := inject(ctx, reflect.ValueOf(object), src); err != nil {
		return nil, err
	}
	return object, nil
}

func classify(ctx context.Context, src reflect.Value) (string, error) {
	if src.IsValid() && src.Type() == types.Any {
		src = src.Elem()
	}
	if src.Kind() != reflect.Map {
		return "", fmt.Errorf("expect map but %v", src.Kind())
	}
	key := ctxGetObjectKey(ctx)
	rkind := src.MapIndex(reflect.ValueOf(key))
	if !rkind.IsValid() {
		return "", fmt.Errorf("not exist kind field(%s)", key)
	}
	if rkind.Type() == types.Any {
		rkind = rkind.Elem()
	}
	if rkind.Type() != types.String {
		return "", fmt.Errorf("kind must be a string")
	}
	return rkind.Interface().(string), nil
}
