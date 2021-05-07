package tracking

import (
	"context"
	"fmt"

	"github.com/blacklane/go-libs/x/events"
	"github.com/opentracing/opentracing-go"
)

// SpanFromContext returns the non-nil span returned by opentracing.SpanFromContext
// or a opentracing.noopSpan if there is no span in the context.
func SpanFromContext(ctx context.Context) opentracing.Span {
	sp := opentracing.SpanFromContext(ctx)
	if sp == nil {
		return opentracing.NoopTracer{}.StartSpan("noopTracer")
	}

	return sp
}

// EventsOpentracingInject injects span into event.Headers. It's safe to call
// this function if event.Headers is nil.
// A non-nil error is returned on failure.
func EventsOpentracingInject(tracer opentracing.Tracer, span opentracing.Span, event *events.Event) error {
	if event.Headers == nil {
		event.Headers = map[string]string{}
	}

	err := tracer.Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(event.Headers))
	if err != nil {
		return fmt.Errorf("could not inject opentracing span: %w", err)
	}

	return nil
}
