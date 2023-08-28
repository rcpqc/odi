package resolve

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/container"
	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

// IResolve custom resolve interface for an object
type IResolve interface{ Resolve(src any) error }

func invoke(ctx context.Context, src reflect.Value) (any, error) {
	kind, err := classify(src, ctxKey(ctx))
	if err != nil {
		return nil, errs.Newf("classify err: %v", err)
	}
	object, err := container.Create(kind)
	if err != nil {
		return nil, errs.Newf("container create err: %v", err)
	}
	if iface, ok := object.(IResolve); ok {
		if err := iface.Resolve(src.Interface()); err != nil {
			return nil, err
		}
	} else {
		if err := inject(ctx, reflect.ValueOf(object), src); err != nil {
			return nil, err
		}
	}
	return object, nil
}

func classify(src reflect.Value, key string) (string, error) {
	if src.Kind() != reflect.Map || src.Type().Key() != types.String {
		return "", fmt.Errorf("expected map[string]any")
	}
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
