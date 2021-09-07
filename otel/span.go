package otel

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SpanAddErr adds an error to span with SpanErrorMessageKey and SpanErrorTypeKey
// attributes.
func SpanAddErr(span trace.Span, err error) {
	span.RecordError(err, trace.WithStackTrace(true))
	span.SetStatus(codes.Error, err.Error())
}
