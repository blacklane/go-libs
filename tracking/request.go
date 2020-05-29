package tracking

import "context"

type (
	depthKey        struct{}
	routeKey        struct{}
	treePathContext struct{}
)

var (
	ctxDepthKey    = depthKey{}
	ctxRouteKey    = routeKey{}
	ctxTreePathKey = treePathContext{}
)

func RequestDepthFromCtx(ctx context.Context) int {
	if requestDepth, ok := ctx.Value(ctxDepthKey).(int); ok {
		return requestDepth
	}
	return 0
}

func WithRequestDepth(ctx context.Context, requestDepth int) context.Context {
	return context.WithValue(ctx, ctxDepthKey, requestDepth)
}

// RequestRouteFromCtx returns the request rout in the context or empty string if not found
func RequestRouteFromCtx(ctx context.Context) string {
	if requestRoute, ok := ctx.Value(ctxRouteKey).(string); ok {
		return requestRoute
	}
	return ""
}

// WithRequestRoute TODO
func WithRequestRoute(ctx context.Context, requestRoute string) context.Context {
	return context.WithValue(ctx, ctxRouteKey, requestRoute)
}

// TreePathFromCtx returns the tree path in the context or empty string if not found
func TreePathFromCtx(ctx context.Context) string {
	if treePath, ok := ctx.Value(ctxTreePathKey).(string); ok {
		return treePath
	}
	return ""
}

func WithTreePath(ctx context.Context, treePath string) context.Context {
	return context.WithValue(ctx, ctxTreePathKey, treePath)
}
