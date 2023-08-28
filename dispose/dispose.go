package dispose

import (
	"context"
	"reflect"
)

// IDispose custom dispose interface for an object
type IDispose interface{ Dispose() error }

type Option func(ctx context.Context) context.Context

func WithTagKey(key string) Option {
	return func(ctx context.Context) context.Context { return ctxWithTagKey(ctx, key) }
}

func Destory(target any, opts ...Option) error {
	ctx := context.Background()
	for _, opt := range opts {
		ctx = opt(ctx)
	}
	return destory(ctx, reflect.ValueOf(target))
}
