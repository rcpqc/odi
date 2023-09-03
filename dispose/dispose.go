package dispose

import (
	"context"
	"reflect"
)

// IDispose custom dispose interface for an object
type IDispose interface{ Dispose() error }

type Option func(ctx context.Context) context.Context

func Destory(target any, opts ...Option) error {
	ctx := context.Background()
	return destory(ctx, reflect.ValueOf(target))
}
