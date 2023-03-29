package otel

import (
	"context"
	"net/http"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/otel/internal/constants"
	"github.com/blacklane/go-libs/tracking"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

type responseWriterWrapper struct {
	w    http.ResponseWriter
	span trace.Span
}

func (rwh *responseWriterWrapper) Header() http.Header {
	return rwh.w.Header()
}

func (rwh *responseWriterWrapper) Write(bytes []byte) (int, error) {
	return rwh.w.Write(bytes)
}

func (rwh *responseWriterWrapper) WriteHeader(statusCode int) {
	rwh.span.SetAttributes(semconv.HTTPStatusCodeKey.Int(statusCode))
	rwh.w.WriteHeader(statusCode)
}

// HTTPMiddleware returns a http.NewHandler and adds HTTP information on the span (similar to
// the otelhttp.NewHandler middleware) and adds extra information on top.
// Use go.opentelemetry.io/otel/trace.SpanFromContext to get the span from the context.
func HTTPMiddleware(serviceName, handlerName, path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			propagator := otel.GetTextMapPropagator()
			ctx := r.Context()
			ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

			trackingID := tracking.IDFromContext(ctx)

			attributes := []attribute.KeyValue{
				semconv.HTTPRouteKey.String(path),
				AttrKeyTrackingID.String(trackingID),
			}
			attributes = append(attributes, semconv.NetAttributesFromHTTPRequest("tcp", r)...)
			attributes = append(attributes, semconv.EndUserAttributesFromHTTPRequest(r)...)
			attributes = append(attributes, semconv.HTTPServerAttributesFromHTTPRequest(serviceName, "", r)...)

			ctx, span := otel.Tracer(constants.TracerName).Start(
				ctx,
				handlerName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(attributes...),
			)
			defer span.End()

			log := logger.From(ctx)
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str(constants.LogKeyTraceID, span.SpanContext().TraceID().String())
			})
			ctx = log.WithContext(ctx)

			ww := &responseWriterWrapper{w, span}
			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}

// NewHTTPHandler wraps the passed handler in a span named like operation.
func NewHTTPHandler(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sp := trace.SpanFromContext(ctx)

			trackingID := tracking.IDFromContext(ctx)
			sp.SetAttributes(
				AttrKeyTrackingID.String(trackingID),
			)

			traceID := sp.SpanContext().TraceID()
			log := logger.From(ctx).
				With().
				Stringer(constants.LogKeyTraceID, traceID).
				Logger()

			ctx = log.WithContext(ctx)

			handler.ServeHTTP(w, r.WithContext(ctx))
		},
	), operation)
}

// HTTPInject injects OTel "cross-cutting concerns" (a.k.a OTel headers) and
// X-Tracking-Id into the outgoing request headers.
func HTTPInject(ctx context.Context, r *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
	r.Header.Set("X-Tracking-Id", tracking.IDFromContext(ctx))
}

type otelTransport struct {
	Transport http.RoundTripper
}

func (o *otelTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	HTTPInject(ctx, req)
	return o.Transport.RoundTrip(req)
}

// InstrumentTransport instruments the given http.RoundTripper with OTel.
// If the given http.RoundTripper is nil, http.DefaultTransport is used.
// Requests should be made with the context with the span.
func InstrumentTransport(t http.RoundTripper) http.RoundTripper {
	if t == nil {
		t = http.DefaultTransport
	}
	return &otelTransport{Transport: t}
}
