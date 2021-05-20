package tracing

import (
	"context"
	"testing"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/opentracing/opentracing-go"

	"github.com/blacklane/go-libs/tracing/internal/jeager"
)

func TestEventsAddOpentracing_noSpan(t *testing.T) {
	wantTrackingID := "TestEventsAddOpentracing"
	handler := events.HandlerFunc(func(ctx context.Context, e events.Event) error {
		validateRootSpan(t, SpanFromContext(ctx), wantTrackingID)
		return nil
	})

	tracer, closer := jeager.NewTracer("", "Opentracing-integration-test", noopLogger)
	defer func() {
		if err := closer.Close(); err != nil {
			t.Errorf("error closing tracer: %v", err)
		}
	}()

	h := EventsAddOpentracing("TestEventsAddOpentracing", tracer)(handler)
	ctx := tracking.SetContextID(context.Background(), wantTrackingID)
	_ = h.Handle(ctx, events.Event{})
}

func TestEventsAddOpentracing_existingSpan(t *testing.T) {
	e := events.Event{Headers: map[string]string{}}
	eventName := "TestEventsAddOpentracing_existingSpan"
	wantTrackingID := "TestEventsAddOpentracing"
	handler := events.HandlerFunc(func(ctx context.Context, e events.Event) error {
		validateChildSpan(t, SpanFromContext(ctx), wantTrackingID)
		return nil
	})

	tracer, closer := jeager.NewTracer("", "Opentracing-integration-test", noopLogger)
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

	h := EventsAddOpentracing("TestEventsAddOpentracing", tracer)(handler)
	_ = h.Handle(tracking.SetContextID(context.Background(), wantTrackingID), e)
}
