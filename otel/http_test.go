package otel

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/blacklane/go-libs/tracking"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	stdoutExporter, err := stdouttrace.New()
	if err != nil {
		panic(fmt.Sprintf(
			"failed to initialize stdouttrace export pipeline: %v", err))
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(stdoutExporter),
	)

	tracerProvider.RegisterSpanProcessor(
		sdktrace.NewSimpleSpanProcessor(stdoutExporter))

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func TestHTTPInject(t *testing.T) {
	const spanName = "TestHTTPInject"
	const trackingID = "trackingID_" + spanName

	ctx, _ := otel.GetTracerProvider().Tracer(spanName).
		Start(context.Background(), spanName)
	ctx = tracking.SetContextID(ctx, trackingID)

	r, err := http.NewRequest(http.MethodGet, "/ignore", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	HTTPInject(ctx, r)

	gotTrackingID := r.Header.Get("X-Tracking-Id")
	if gotTrackingID != trackingID {
		t.Errorf("got X-Tracking-Id = %s, want %s", gotTrackingID, trackingID)
	}

	gotCtx := otel.GetTextMapPropagator().
		Extract(context.Background(), propagation.HeaderCarrier(r.Header))
	if !trace.SpanFromContext(gotCtx).SpanContext().IsValid() {
		t.Errorf("got a span with invalid span context, want a valid span context")
	}
}
