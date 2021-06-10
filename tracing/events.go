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
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blacklane/go-libs/tracing/internal/constants"
)

const (
	OtelKeyEventName  = attribute.Key("event.name")
	OtelKeyTrackingID = attribute.Key("tracking_id")
)

// EventsAddDefault returns the composition of EventsGenericMiddleware and
// EventsHandlerMiddleware. Therefore it should be applied to each events.Handler
// individually.
func EventsAddDefault(handler events.Handler, log logger.Logger, eventName string) events.Handler {
	hb := events.HandlerBuilder{}
	hb.AddHandler(handler)
	hb.UseMiddleware(
		EventsGenericMiddleware(log),
		EventsHandlerMiddleware(eventName),
	)

	return hb.Build()[0]
}

// EventsGenericMiddleware adds all middlewares which aren't specific to a handler.
// They are:
// - github.com/blacklane/go-libs/tracking/middleware.EventsAddTrackingID
// - github.com/blacklane/go-libs/logger/middleware.EventsAddLogger
// It adds log as the logger and will not log any route in skipRoutes.
func EventsGenericMiddleware(log logger.Logger) events.Middleware {
	return func(handler events.Handler) events.Handler {
		hb := events.HandlerBuilder{}
		hb.AddHandler(handler)
		hb.UseMiddleware(
			trackmiddleware.EventsAddTrackingID,
			logmiddleware.EventsAddLogger(log),
		)

		return hb.Build()[0]
	}
}

// EventsHandlerMiddleware are event specific middleware. It adds:
// - EventsAddOpenTelemetry
// - github.com/blacklane/go-libs/logger/middleware.EventsHandlerStatusLogger
func EventsHandlerMiddleware(eventName string) events.Middleware {
	return func(handler events.Handler) events.Handler {
		hb := events.HandlerBuilder{}
		hb.AddHandler(handler)
		hb.UseMiddleware(
			EventsAddOpenTelemetry(eventName),
			logmiddleware.EventsHandlerStatusLogger(eventName),
		)

		return hb.Build()[0]
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
			// Get the global TracerProvider. // TODO:make it a constant and have one for http
			tr := otel.Tracer(constants.TracerName)

			trackingID := eventsExtractTrackingID(ctx, e)

			ctx = propagator.Extract(ctx, e.Headers)

			ctx, sp := tr.Start(
				ctx,
				eventName,
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					OtelKeyTrackingID.String(trackingID),
					OtelKeyEventName.String(eventName)),
			)
			defer sp.End()

			logger.FromContext(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str(constants.TraceID, sp.SpanContext().TraceID().String())
			})

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
