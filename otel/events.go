package otel

import (
	"context"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/blacklane/go-libs/x/events/consumer"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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
// Use go.opentelemetry.io/otel/trace.SpanFromContext to get the span from the context.
func EventsAddOpenTelemetry(eventName string) events.Middleware {
	return func(handler events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
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
				return c.Str(constants.LogKeyTraceID, sp.SpanContext().TraceID().String())
			})

			return handler.Handle(ctx, e)
		})
	}
}

func EventConsumer(eventsToInclude ...string) consumer.Middleware {
	return func(next consumer.Handler) consumer.Handler {
		return func(ctx context.Context, m consumer.Message) error {
			eventName := m.EventName()

			if !isEventIncluded(eventName, eventsToInclude) {
				return next(ctx, m)
			}

			tr := otel.Tracer(constants.TracerName)
			trackingID := tracking.IDFromContext(ctx)

			ctx = otel.GetTextMapPropagator().Extract(ctx, m.Header())

			ctx, sp := tr.Start(
				ctx,
				eventName,
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					AttrKeyTrackingID.String(trackingID),
					AttrKeyEventName.String(eventName),
					semconv.MessagingMessageIDKey.String(m.ID()),
					semconv.MessagingKafkaMessageKeyKey.String(m.Key()),
					semconv.MessagingKafkaPartitionKey.Int(int(m.TopicPartition().Partition)),
				),
			)
			defer sp.End()

			log := logger.FromContext(ctx).
				With().
				Str(constants.LogKeyEventName, eventName).
				Stringer(constants.LogKeyTraceID, sp.SpanContext().TraceID()).
				Logger()

			ctx = log.WithContext(ctx)

			if err := next(ctx, m); err != nil {
				SpanAddErr(sp, err)
				return err
			}

			return nil
		}
	}
}

func isEventIncluded(event string, eventsToInclude []string) bool {
	// We assume all events should be included unless otherwise specified.
	if len(eventsToInclude) == 0 {
		return true
	}

	for _, ev := range eventsToInclude {
		if event == ev {
			return true
		}
	}

	return false
}
