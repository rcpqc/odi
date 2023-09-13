package resolve

import (
	"context"
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
	return invoke(ctx, reflect.ValueOf(src))
}

func Inject(dst, src any, opts ...Option) error {
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	return injectIgnoreResolve(ctx, reflect.ValueOf(dst), reflect.ValueOf(src))
}
