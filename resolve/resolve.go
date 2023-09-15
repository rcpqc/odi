package resolve

import (
	"context"
	"fmt"
	"reflect"
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
	rsrc := reflect.ValueOf(src)
	if !rsrc.IsValid() {
		return nil, nil
	}
	return invoke(ctx, rsrc)
}

func Object(dst, src any, opts ...Option) error {
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	rsrc := reflect.ValueOf(src)
	if !rsrc.IsValid() {
		return nil
	}
	return inject(ctx, reflect.ValueOf(dst), rsrc)
}

func Struct(dst, src any, opts ...Option) error {
	rdst := reflect.ValueOf(dst)
	if rdst.Kind() != reflect.Pointer {
		return fmt.Errorf("expect *struct but %v", rdst.Kind())
	}
	if rdst.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expect *struct but *%v", rdst.Elem().Kind())
	}
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	rsrc := reflect.ValueOf(src)
	if !rsrc.IsValid() {
		return nil
	}
	return injectStruct(ctx, rdst.Elem(), rsrc)
}
