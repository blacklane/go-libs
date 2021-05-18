package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"
	"github.com/opentracing/opentracing-go"

	"github.com/blacklane/go-libs/tracking"
)

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

// EventsAddOpentracing adds an opentracing span to the context and finishes the span
// when the handler returns.
// Use tracking.SpanFromContext to get the span from the context. It is
// technically safe to call opentracing.SpanFromContext after this middleware
// and trust the returned span is not nil. However tracking.SpanFromContext is
// safer as it'll return a disabled span if none is found in the context.
func EventsAddOpentracing(eventName string, tracer opentracing.Tracer, handler events.Handler) events.Handler {
	return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
		trackingID := eventsExtractTrackingID(ctx, e)

		span := eventGetChildSpan(ctx, e, eventName, tracer)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)

		// Set as not all systems will update to opentracing at once
		span.SetTag("tracking_id", trackingID)

		if e.Headers == nil {
			e.Headers = map[string]string{}
		}

		return handler.Handle(ctx, e)
	})
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

func eventGetChildSpan(ctx context.Context, e events.Event, eName string, tracer opentracing.Tracer) opentracing.Span {
	carrier := opentracing.TextMapCarrier(e.Headers)

	spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
	if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
		logger.FromContext(ctx).
			Err(err).
			Msg("tracing events: could not extract span")
	}

	span := tracer.StartSpan(eName, opentracing.ChildOf(spanContext))

	return span
}
