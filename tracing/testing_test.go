package tracing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/blacklane/go-libs/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	"github.com/blacklane/go-libs/tracing/internal/constants"
)

var noopLogger = *logger.FromContext(context.Background())

func validateChildSpan(t *testing.T, span opentracing.Span, wantTrackingID string) {
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

func validateRootSpan(t *testing.T, span opentracing.Span, wantTrackingID string) {
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

// prettyJSONWriter takes a JSON string as []byte an ident it before writing to os.Stdout..
type prettyJSONWriter struct{}

func (pw prettyJSONWriter) Write(p []byte) (int, error) {
	buf := &bytes.Buffer{}

	data := &map[string]interface{}{}
	if err := json.Unmarshal(p, data); err != nil {
		return 0, fmt.Errorf("prettyJSONWriter: could not json.Unmarshal data: %w", err)
	}

	pp, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("prettyJSONWriter: could not json.MarshalIndent data: %w", err)
	}

	buf.Write(pp)
	buf.WriteByte('\n')

	n, err := buf.WriteTo(os.Stdout)
	return int(n), err
}
