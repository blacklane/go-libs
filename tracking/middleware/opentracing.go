package middleware

import (
	"context"
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

		span := httpExtractSpan(r, path, tracer)
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
		span := eventExtractSpan(ctx, e, eventName, tracer)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)

		// Set as not all systems will update to opentracing at once
		span.SetTag("tracking_id", tracking.IDFromContext(ctx))

		err := tracer.Inject(
			span.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(e.Headers))
		if err != nil {
			logger.FromContext(ctx).Err(err).Msg("could not inject opentracing span")
		}

		return handler.Handle(ctx, e)
	})
}

func httpExtractSpan(r *http.Request, path string, tracer opentracing.Tracer) opentracing.Span {
	carrier := opentracing.HTTPHeadersCarrier(r.Header)

	spanContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err != nil {
		logger.FromContext(r.Context()).Err(err).Msg("could not extract span")
	}

	span := tracer.StartSpan(
		path,
		opentracing.ChildOf(spanContext))

	return span
}

func eventExtractSpan(ctx context.Context, ev events.Event, evName string, tracer opentracing.Tracer) opentracing.Span {
	carrier := opentracing.TextMapCarrier(ev.Headers)

	spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
	if err != nil {
		logger.FromContext(ctx).Err(err).Msg("could not extract span")
	}

	span := tracer.StartSpan(
		evName,
		opentracing.ChildOf(spanContext))

	return span
}
