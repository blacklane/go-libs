package middleware

import (
	"context"

	"github.com/blacklane/go-libs/x/events"

	"github.com/blacklane/go-libs/logger"
)

// EventsAddLogger add the logger into the context
func EventsAddLogger(log logger.Logger) events.Middleware {
	return func(next events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			return next.Handle(log.WithContext(ctx), e)
		})
	}
}
