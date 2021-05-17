package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/tracking/internal/constants"
	"github.com/blacklane/go-libs/tracking/internal/jeager"
)

var noopLogger = *logger.FromContext(context.Background())

func TestHTTPAddOpentracing_noSpan(t *testing.T) {
	wantTrackingID := "TestHTTPAddOpentracing_noSpan"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateRootSpan(t, r.Context(), wantTrackingID)
	})

	tracer, closer := jeager.NewTracer("Opentracing-integration-test", noopLogger)
	defer func() {
		if err := closer.Close(); err != nil {
			t.Errorf("error closing tracer: %v", err)
		}
	}()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(tracking.SetContextID(r.Context(), wantTrackingID))

	_, h := HTTPAddOpentracing("/TestHTTPAddOpentracing_noSpan", tracer, handler)
	h.ServeHTTP(httptest.NewRecorder(), r)
}

func TestHTTPAddOpentracing_existingSpan(t *testing.T) {
	wantTrackingID := "TestHTTPAddOpentracing_noSpan"
	path := "/TestHTTPAddOpentracing_noSpan"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateChildSpan(t, r.Context(), wantTrackingID)
	})

	tracer, closer := jeager.NewTracer("Opentracing-integration-test", noopLogger)
	defer func() {
		if err := closer.Close(); err != nil {
			t.Errorf("error closing tracer: %v", err)
		}
	}()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(tracking.SetContextID(r.Context(), wantTrackingID))

	span := tracer.StartSpan(path)
	err := tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		t.Fatalf("could not inject span into event headers: %v", err)
	}

	_, h := HTTPAddOpentracing(path, tracer, handler)
	h.ServeHTTP(httptest.NewRecorder(), r)
}

func TestEventsAddOpentracing_noSpan(t *testing.T) {
	wantTrackingID := "TestEventsAddOpentracing"
	handler := events.HandlerFunc(func(ctx context.Context, e events.Event) error {
		validateRootSpan(t, ctx, wantTrackingID)
		return nil
	})

	tracer, closer := jeager.NewTracer("Opentracing-integration-test", noopLogger)
	defer func() {
		if err := closer.Close(); err != nil {
			t.Errorf("error closing tracer: %v", err)
		}
	}()

	h := EventsAddOpentracing("TestEventsAddOpentracing", tracer, handler)
	_ = h.Handle(tracking.SetContextID(context.Background(), wantTrackingID), events.Event{})
}

func TestEventsAddOpentracing_existingSpan(t *testing.T) {
	e := events.Event{Headers: map[string]string{}}
	eventName := "TestEventsAddOpentracing_existingSpan"
	wantTrackingID := "TestEventsAddOpentracing"
	handler := events.HandlerFunc(func(ctx context.Context, e events.Event) error {
		validateChildSpan(t, ctx, wantTrackingID)
		return nil
	})

	tracer, closer := jeager.NewTracer("Opentracing-integration-test", noopLogger)
	defer func() {
		if err := closer.Close(); err != nil {
			t.Errorf("error closing tracer: %v", err)
		}
	}()

	span := tracer.StartSpan(eventName)
	err := tracer.Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(e.Headers))
	if err != nil {
		t.Fatalf("could not inject span into event headers: %v", err)
	}

	h := EventsAddOpentracing("TestEventsAddOpentracing", tracer, handler)
	_ = h.Handle(tracking.SetContextID(context.Background(), wantTrackingID), e)
}

func validateChildSpan(t *testing.T, ctx context.Context, wantTrackingID string) {
	span := tracking.SpanFromContext(ctx)

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
	span := tracking.SpanFromContext(ctx)

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
