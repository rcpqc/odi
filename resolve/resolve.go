package resolve

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/types"
)

type Option func(ctx context.Context) context.Context

func WithObjectKey(key string) Option {
	return func(ctx context.Context) context.Context { return ctxWithObjectKey(ctx, key) }
}

func WithTagKey(key string) Option {
	return func(ctx context.Context) context.Context { return ctxWithTagKey(ctx, key) }
}

func WithStructFieldNameCompatibility(enable bool) Option {
	return func(ctx context.Context) context.Context { return ctxWithStructFieldNameCompatibility(ctx, enable) }
}

func Invoke(src any, opts ...Option) (any, error) {
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	return invoke(ctx, reflect.ValueOf(src))
}

func Struct(dst, src any, opts ...Option) error {
	rdst := reflect.ValueOf(dst)
	if rdst.Kind() != reflect.Pointer || rdst.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expect *struct but %v", rdst.Kind())
	}
	if rdst.IsNil() {
		return fmt.Errorf("*struct is nil")
	}
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	rsrc := reflect.ValueOf(src)
	if rsrc.IsValid() && rsrc.Type() == types.Any {
		rsrc = rsrc.Elem()
	}
	return injectStruct(ctx, rdst.Elem(), rsrc)
}
