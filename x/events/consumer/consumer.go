package consumer

import (
	"context"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"
)

type consumer struct {
	eventName string
	handler   Handler
}

func New(eventName string, handler Handler, middlewares ...Middleware) events.Handler {
	return &consumer{
		eventName: eventName,
		handler:   ApplyMiddlewares(handler, middlewares),
	}
}

func (l *consumer) Handle(ctx context.Context, e events.Event) error {
	log := *logger.FromContext(ctx)
	m, err := createJsonMessage(e)
	if err != nil {
		log.Err(err).Msgf("could not unmarshal kafka event: %s", e.Payload)
		return err
	}

	if m.EventName() != l.eventName {
		return nil
	}

	log = log.With().
		Str("event_name", l.eventName).
		Logger()

	ctx = log.WithContext(ctx)
	return l.handler(ctx, m)
}
