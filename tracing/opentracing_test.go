package tracing

import (
	"context"
	"testing"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"
	"github.com/opentracing/opentracing-go"

	"github.com/blacklane/go-libs/tracing/internal/jeager"
)

func TestSpanFromContext(t *testing.T) {
	noopSpan := opentracing.NoopTracer{}.StartSpan("noopTracer")
	tracer, closer :=
		jeager.NewTracer("Opentracing-test", *logger.FromContext(context.Background()))
	defer closer.Close()

	want := tracer.StartSpan("TestSpanFromContext")
	ctx := opentracing.ContextWithSpan(
		context.Background(), want)

	got := SpanFromContext(ctx)

	if got == noopSpan {
		t.Errorf("got a noop span, want: %T", want)
	}
}

func TestSpanFromContext_noopSpan(t *testing.T) {
	want := opentracing.NoopTracer{}.StartSpan("noopTracer")

	got := SpanFromContext(context.Background())

	if got != want {
		t.Errorf("got: %T, want: %T", got, want)
	}
}

func TestEventsOpentracingInject(t *testing.T) {
	tracer, closer :=
		jeager.NewTracer("Opentracing-test", *logger.FromContext(context.Background()))
	defer closer.Close()

	e := events.Event{}
	span := tracer.StartSpan("TestEventsOpentracingInject")
	err := EventsOpentracingInject(context.Background(), span, &e)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	_, err = tracer.Extract(
		opentracing.TextMap, opentracing.TextMapCarrier(e.Headers))
	if err != nil {
		t.Errorf("want nil error, got: %v", err)
	}
}
