package tracing

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// OTel span attribute keys.
const (
	SpanErrorMessageKey = attribute.Key("error.message")
	SpanErrorTypeKey    = attribute.Key("error.type")
)

// SpanAddErr adds an error to span with SpanErrorMessageKey and SpanErrorTypeKey
// attributes.
func SpanAddErr(span trace.Span, err error) {
	// TODO(Anderson): Use/test span.RecordError() instead of adding it manually
	span.AddEvent("error", trace.WithAttributes(
		SpanErrorMessageKey.String(err.Error()),
		SpanErrorTypeKey.String(fmt.Sprintf("%T", err)),
	))
}
