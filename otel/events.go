package otel

import (
	"context"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blacklane/go-libs/otel/internal/constants"
)

// OpenTelemetry attribute keys.
const (
	AttrKeyEventName  = attribute.Key("event.name")
	AttrKeyTrackingID = attribute.Key("tracking.id")
)

// EventsAddOpenTelemetry adds an opentracing span to the context and finishes the span
// when the handler returns.
// Use tracking.SpanFromContext to get the span from the context. It is
// technically safe to call opentracing.SpanFromContext after this middleware
// and trust the returned span is not nil. However tracking.SpanFromContext is
// safer as it'll return a disabled span if none is found in the context.
func EventsAddOpenTelemetry(eventName string) events.Middleware {
	return func(handler events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			// Get the global TracerProvider. // TODO:make it a constant and have one for http
			tr := otel.Tracer(constants.TracerName)

			trackingID := tracking.IDFromContext(ctx)

			ctx = otel.GetTextMapPropagator().Extract(ctx, e.Headers)

			ctx, sp := tr.Start(
				ctx,
				eventName,
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					AttrKeyTrackingID.String(trackingID),
					AttrKeyEventName.String(eventName)),
			)
			defer sp.End()

			logger.FromContext(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str(constants.TraceID, sp.SpanContext().TraceID().String())
			})

			return handler.Handle(ctx, e)
		})
	}
}
