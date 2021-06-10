package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

// SpanAddErr adds an error event to span with error.message and error.type
// attributes.
func SpanAddErr(span trace.Span, err error) {
	span.AddEvent("error", trace.WithAttributes(
		attribute.String("error.message", err.Error()),
		attribute.String("error.type", fmt.Sprintf("%T", err)),
	))
}

// SpanStartProducedEvent starts a span with tracerName=kind=producer and
// MessagingDestinationKey=topic.
// The caller must end the span, usually deferring it:
//    defer span.End()
func SpanStartProducedEvent(ctx context.Context, eventName string, topic string) (context.Context, trace.Span) {
	return trace.SpanFromContext(ctx).Tracer().
		Start(ctx, "producer:"+eventName,
			trace.WithSpanKind(trace.SpanKindProducer),
			trace.WithAttributes(semconv.MessagingDestinationKey.String(topic)))
}
