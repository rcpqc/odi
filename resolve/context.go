package resolve

import "context"

var (
	ctxDefaultObjectKey = "object"
	ctxDefaultTagKey    = "odi"
)

type ctxObject struct{}

func ctxWithObjectKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, ctxObject{}, key)
}
func ctxGetObjectKey(ctx context.Context) string {
	val, ok := ctx.Value(ctxObject{}).(string)
	if !ok {
		return ctxDefaultObjectKey
	}
	return val
}

type ctxTag struct{}

func ctxWithTagKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, ctxTag{}, key)
}
func ctxGetTagKey(ctx context.Context) string {
	val, ok := ctx.Value(ctxTag{}).(string)
	if !ok {
		return ctxDefaultTagKey
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
