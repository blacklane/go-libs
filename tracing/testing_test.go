package tracing

import (
	"context"
	"testing"

	"github.com/blacklane/go-libs/logger"
	"github.com/uber/jaeger-client-go"

	"github.com/blacklane/go-libs/tracing/internal/constants"
)

var noopLogger = *logger.FromContext(context.Background())

func validateChildSpan(t *testing.T, ctx context.Context, wantTrackingID string) {
	span := SpanFromContext(ctx)

	jeagerSpan, ok := span.(*jaeger.Span)
	if !ok {
		t.Errorf("want span of type %T, got: %T", &jaeger.Span{}, span)
	}
	if trackingID := jeagerSpan.Tags()[constants.TagTrackingID]; trackingID == "" {
		t.Errorf("want trackingID: %s, got: %s", wantTrackingID, trackingID)
	}
	if jeagerSpan.SpanContext().ParentID() == 0 {
		t.Errorf("want a child span, got a root span")
	}
}

func validateRootSpan(t *testing.T, ctx context.Context, wantTrackingID string) {
	span := SpanFromContext(ctx)

	jeagerSpan, ok := span.(*jaeger.Span)
	if !ok {
		t.Errorf("want span of type %T, got: %T", &jaeger.Span{}, span)
	}
	if trackingID := jeagerSpan.Tags()[constants.TagTrackingID]; trackingID == "" {
		t.Errorf("want trackingID: %s, got: %s", wantTrackingID, trackingID)
	}
	if jeagerSpan.SpanContext().ParentID() != 0 {
		t.Errorf("want a root span, got a child span")
	}
}
