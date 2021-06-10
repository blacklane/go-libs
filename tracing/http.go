package tracing

import (
	"net/http"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	"github.com/blacklane/go-libs/tracking"
	trackingmiddleware "github.com/blacklane/go-libs/tracking/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"

	"github.com/blacklane/go-libs/tracing/internal/constants"
)

// HTTPAllMiddleware returns the composition of HTTPGenericMiddleware and
// HTTPRequestMiddleware. Therefore it should be applied to each http.Handler
// individually.
func HTTPAllMiddleware(serviceName, handlerName, path string, log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := HTTPRequestMiddleware(serviceName, handlerName, path)(next)
			h = HTTPGenericMiddleware(log, skipRoutes...)(h)

			h.ServeHTTP(w, r)
		})
	}
}

// HTTPGenericMiddleware adds all middlewares which aren't specific to a handler.
// They are:
// - github.com/blacklane/go-libs/tracking/middleware.TrackingID
// - github.com/blacklane/go-libs/logger/middleware.HTTPAddLogger
// - github.com/blacklane/go-libs/logger/middleware.HTTPRequestLogger
// It adds log as the logger and will not log any route in skipRoutes.
func HTTPGenericMiddleware(log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := logmiddleware.HTTPRequestLogger(skipRoutes)(next)
			h = logmiddleware.HTTPAddLogger(log)(h)
			h = trackingmiddleware.TrackingID(h)
			h.ServeHTTP(w, r)
		})
	}
}

// HTTPRequestMiddleware are request specific middleware. It adds OpenTelemetry
// otelhttp.NewHandler and extra information on the span created by the OTel
// middleware.
func HTTPRequestMiddleware(serviceName, handlerName, path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			trackingID := tracking.IDFromContext(ctx)

			sp := trace.SpanFromContext(ctx)
			sp.SetName(handlerName)
			sp.SetAttributes(
				semconv.HTTPRouteKey.String(path), // same as adding otelhttp.WithRouteTag
				OtelKeyTrackingID.String(trackingID))

			logger.FromContext(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str(constants.TraceID, sp.SpanContext().TraceID().String())
			})

			next.ServeHTTP(w, r)
		}), serviceName)
	}
}
