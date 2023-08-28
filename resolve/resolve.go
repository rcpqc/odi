package resolve

import (
	"context"
	"reflect"
)

type Option func(ctx context.Context) context.Context

func WithKey(key string) Option {
	return func(ctx context.Context) context.Context { return ctxWithKey(ctx, key) }
}

func Invoke(src any, opts ...Option) (any, error) {
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	return invoke(ctx, reflect.ValueOf(src))
}

func Inject(dst, src any, opts ...Option) error {
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	return inject(ctx, reflect.ValueOf(dst), reflect.ValueOf(src))
}
