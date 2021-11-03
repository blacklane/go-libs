package otel

import (
	"context"
	"net/http"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/blacklane/go-libs/otel/internal/constants"
)

// HTTPMiddleware adds OpenTelemetry otelhttp.NewHandler and extra information
// on the span created by the otelhttp.NewHandler middleware.
// Use go.opentelemetry.io/otel/trace.SpanFromContext to get the span from the context.
func HTTPMiddleware(serviceName, handlerName, path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				trackingID := tracking.IDFromContext(ctx)

				sp := trace.SpanFromContext(ctx)
				sp.SetName(handlerName)
				sp.SetAttributes(
					semconv.HTTPRouteKey.String(path), // same as adding otelhttp.WithRouteTag
					AttrKeyTrackingID.String(trackingID))

				logger.FromContext(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str(constants.LogKeyTraceID, sp.SpanContext().TraceID().String())
				})

				next.ServeHTTP(w, r)
			}),
			serviceName)
	}
}

// HTTPInject injects OTel "cross-cutting concerns" (a.k.a OTel headers) and
// X-Tracking-Id into the outgoing request headers.
func HTTPInject(ctx context.Context, r *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
	r.Header.Set("X-Tracking-Id", tracking.IDFromContext(ctx))
}
