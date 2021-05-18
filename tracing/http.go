package tracing

import (
	"errors"
	"net/http"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	"github.com/blacklane/go-libs/tracking"
	trackingmiddleware "github.com/blacklane/go-libs/tracking/middleware"
	"github.com/opentracing/opentracing-go"
)

// HTTPAddDefault adds the necessary middleware for:
//   - have tracking id in the context (read from the headers or a new one),
//   - have a logger.Logger with tracking id and all required fields in the context,
//   - log, at the end of handler, if it succeeded or failed and how log it took.
// - github.com/blacklane/go-libs/tracking/middleware.TrackingID
// - middleware.HTTPAddLogger
// - middleware.HTTPRequestLogger
func HTTPAddDefault(log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := logmiddleware.HTTPRequestLogger(skipRoutes)(next)
			h = logmiddleware.HTTPAddLogger(log)(h)
			h = trackingmiddleware.TrackingID(h)
			h.ServeHTTP(w, r)
		})
	}
}

// HTTPAddOpentracing adds an opentracing span to the context and finishes the span
// when the handler returns.
// Use tracking.SpanFromContext to get the span from the context. It is
// technically safe to call opentracing.SpanFromContext after this middleware
// and trust the returned span is not nil. However tracking.SpanFromContext is
// safer as it'll return a disabled span if none is found in the context.
func HTTPAddOpentracing(path string, tracer opentracing.Tracer, handler http.Handler) (string, http.HandlerFunc) {
	return path, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := httpGetChildSpan(r, path, tracer)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)

		// Set as not all systems will update to opentracing at once
		span.SetTag("tracking_id", tracking.IDFromContext(ctx))

		err := tracer.Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(w.Header()))
		if err != nil {
			logger.FromContext(ctx).Err(err).Msg("could not inject opentracing span")
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

func httpGetChildSpan(r *http.Request, path string, tracer opentracing.Tracer) opentracing.Span {
	carrier := opentracing.HTTPHeadersCarrier(r.Header)

	spanContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
		logger.FromContext(r.Context()).
			Err(err).
			Msg("tracing http: could not extract span")
	}

	span := tracer.StartSpan(path, opentracing.ChildOf(spanContext))

	return span
}
