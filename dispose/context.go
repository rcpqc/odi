package dispose

import "context"

var (
	ctxDefaultTagKey = "odi"
)

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
