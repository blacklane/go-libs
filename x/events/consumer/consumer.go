package consumer

import (
	"context"
	"fmt"

	"github.com/blacklane/go-libs/x/events"
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
		return nil
	}

	return c.handler(ctx, m)
}
