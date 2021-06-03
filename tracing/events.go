package tracing

import (
	"context"
	"errors"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	trackmiddleware "github.com/blacklane/go-libs/tracking/middleware"
	"github.com/blacklane/go-libs/x/events"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blacklane/go-libs/tracing/internal/constants"
)

const KeyTrackingID = attribute.Key("tracking_id")
const KeyError = attribute.Key("error")

// EventsAddDefault adds the necessary middleware for:
//   - have tracking id in the context (read from the headers or a new one),
//   - have a logger.Logger with tracking id and all required fields in the context,
//   - log at the end of handler if it succeeded or failed and how log it took.
// For more details, check the middleware used:
// - github.com/blacklane/go-libs/tracking/middleware.EventsAddTrackingID
// - github.com/blacklane/go-libs/logger/middleware.EventsAddLogger
// - github.com/blacklane/go-libs/logger/middleware.EventsHandlerStatusLogger
// - EventsAddOpentracing
// TODO(Anderson): update docs
func EventsAddDefault(handler events.Handler, log logger.Logger, eventName string) events.Handler {
	hb := events.HandlerBuilder{}
	hb.AddHandler(handler)
	hb.UseMiddleware(
		trackmiddleware.EventsAddTrackingID,
		logmiddleware.EventsAddLogger(log),
		logmiddleware.EventsHandlerStatusLogger(eventName),
		EventsAddOpenTelemetry(eventName))

	return hb.Build()[0]
}

// EventsAddOpentracing adds an opentracing span to the context and finishes the span
// when the handler returns.
// Use tracking.SpanFromContext to get the span from the context. It is
// technically safe to call opentracing.SpanFromContext after this middleware
// and trust the returned span is not nil. However tracking.SpanFromContext is
// safer as it'll return a disabled span if none is found in the context.
func EventsAddOpentracing(eventName string, tracer opentracing.Tracer) events.Middleware {
	return func(handler events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			trackingID := eventsExtractTrackingID(ctx, e)

			span := eventGetChildSpan(ctx, e, eventName, tracer)
			defer span.Finish()

			ctx = opentracing.ContextWithSpan(ctx, span)

			// Set as not all systems will update to opentracing at once
			span.SetTag("tracking_id", trackingID)

			return handler.Handle(ctx, e)
		})
	}
}

// EventsAddOpenTelemetry adds an opentracing span to the context and finishes the span
// when the handler returns.
// Use tracking.SpanFromContext to get the span from the context. It is
// technically safe to call opentracing.SpanFromContext after this middleware
// and trust the returned span is not nil. However tracking.SpanFromContext is
// safer as it'll return a disabled span if none is found in the context.
func EventsAddOpenTelemetry(eventName string) events.Middleware {
	return func(handler events.Handler) events.Handler {
		propagator := otel.GetTextMapPropagator()
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			// Get the global TracerProvider.
			tr := otel.Tracer("github.com/blacklane/go-libs/tracing")

			trackingID := eventsExtractTrackingID(ctx, e)

			ctx = propagator.Extract(ctx, e.Headers)

			ctx, span := tr.Start(ctx,
				eventName,
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(KeyTrackingID.String(trackingID)),
			)
			defer span.End()

			span.AddEvent(eventName)

			// Do bar...
			// defaultOpts := []Option{
			// 	WithSpanOptions(trace.WithSpanKind(trace.SpanKindServer)),
			// 	WithSpanNameFormatter(defaultHandlerFormatter),
			// }
			//
			// opts := append([]trace.SpanOption{
			// 	trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
			// 	trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
			// 	trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(h.operation, "", r)...),
			// }, h.spanStartOptions...) // start with the configured options
			//
			// ctx := h.propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			// ctx, span := h.tracer.Start(ctx, h.spanNameFormatter(h.operation, r), opts...)
			// defer span.End()
			//
			// span := eventGetChildSpan(ctx, e, eventName, tracer)
			// defer span.Finish()
			//
			// ctx = opentracing.ContextWithSpan(ctx, span)
			//
			// // Set as not all systems will update to opentracing at once
			// span.SetTag("tracking_id", trackingID)

			return handler.Handle(ctx, e)
		})
	}
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

func eventsExtractTrackingID(ctx context.Context, e events.Event) string {
	trackingID := e.Headers[constants.HeaderTrackingID]
	if trackingID == "" {
		trackingID = e.Headers[constants.HeaderRequestID]
		if trackingID == "" {
			uuid, err := uuid.NewUUID()
			if err != nil {
				logger.FromContext(ctx).Err(err).Msg("failed to generate trackingID")
				return ""
			}
			trackingID = uuid.String()
		}
	}

	return trackingID
}
