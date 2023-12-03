package resolve

import (
	"context"
	"sync"

	"github.com/rcpqc/odi/types"
)

var (
	ctxDefaultObjectKey      any = "object"
	ctxDefaultProfileFactory     = &types.Factory{TagKey: "odi"}
	ctxProfileFactories          = sync.Map{}
)

type ctxObject struct{}

func ctxWithObjectKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, ctxObject{}, key)
}
func ctxGetObjectKey(ctx context.Context) any {
	val := ctx.Value(ctxObject{})
	if val == nil {
		return ctxDefaultObjectKey
	}
	return val
}

type ctxProfile struct{}

func ctxWithTagKey(ctx context.Context, key string) context.Context {
	val, _ := ctxProfileFactories.LoadOrStore(key, &types.Factory{TagKey: key})
	return context.WithValue(ctx, ctxProfile{}, val)
}
func ctxGetProfileFactory(ctx context.Context) *types.Factory {
	val, ok := ctx.Value(ctxProfile{}).(*types.Factory)
	if !ok {
		return ctxDefaultProfileFactory
	}
	return val
}

type ctxStructFieldNameCompatibility struct{}

func ctxWithStructFieldNameCompatibility(ctx context.Context, enable bool) context.Context {
	return context.WithValue(ctx, ctxStructFieldNameCompatibility{}, enable)
}
func ctxGetStructFieldNameCompatibility(ctx context.Context) bool {
	val, _ := ctx.Value(ctxStructFieldNameCompatibility{}).(bool)
	return val
}
