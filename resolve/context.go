package resolve

import "context"

type ctxKindKey struct{}
type ctxIgnoreKey struct{}

func ctxWithKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, ctxKindKey{}, key)
}
func ctxKey(ctx context.Context) string {
	val, _ := ctx.Value(ctxKindKey{}).(string)
	return val
}

func ctxWithIgnore(ctx context.Context, ignore map[string]struct{}) context.Context {
	return context.WithValue(ctx, ctxIgnoreKey{}, ignore)
}
func ctxIgnore(ctx context.Context) map[string]struct{} {
	val, _ := ctx.Value(ctxIgnoreKey{}).(map[string]struct{})
	return val
}
