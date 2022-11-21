package consumer

import (
	"context"
	"fmt"

	"github.com/blacklane/go-libs/x/events"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ events.Handler = (*Consumer)(nil)

type Consumer struct {
	eventName string
	handler   Handler
}

func New(eventName string, handler Handler, middlewares ...Middleware) *Consumer {
	return &Consumer{
		eventName: eventName,
		handler:   ApplyMiddlewares(handler, middlewares),
	}
}

func (c *Consumer) Handle(ctx context.Context, e events.Event) error {
	m, err := createJsonMessage(e)
	if err != nil {
		return fmt.Errorf("could not unmarshal kafka event: %w", err)
	}

	if m.EventName() != c.eventName {
		if sp := trace.SpanFromContext(ctx); sp.IsRecording() {
			sp.SetAttributes(
				attribute.Bool("event_skipped", true),
			)
		}
		return nil
	}

	return c.handler(ctx, m)
}
