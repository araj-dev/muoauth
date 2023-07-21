package tokenhttp

import (
	"context"
)

type contextIDKey = struct{}

func SetID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextIDKey{}, id)
}

func GetID(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(contextIDKey{}).(string)
	return val, ok
}
