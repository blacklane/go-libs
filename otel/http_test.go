package otel

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestInstrumentTransport(t *testing.T) {
	var traceID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

		remoteSpanCtx := trace.SpanFromContext(ctx).SpanContext()

		if !remoteSpanCtx.IsValid() {
			t.Errorf("got a span with invalid span context")
		}

		remoteTraceID := remoteSpanCtx.TraceID().String()

		if traceID != remoteTraceID {
			t.Errorf("got traceID = %s, want %s", remoteTraceID, traceID)
		}

		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	client := server.Client()
	client.Transport = InstrumentTransport(client.Transport)

	ctx := context.Background()
	ctx, span := otel.Tracer("tracer").Start(ctx, "span")
	defer span.End()

	traceID = span.SpanContext().TraceID().String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("could not create request: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not send request: %s", err)
	}

	if res.StatusCode != http.StatusTeapot {
		t.Errorf("got status code %d, want %d", res.StatusCode, http.StatusTeapot)
	}
}
