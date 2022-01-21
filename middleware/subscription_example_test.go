package middleware_test

import (
	"context"
	"time"

	"github.com/blacklane/go-libs/camunda"
	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/middleware"
)

func ExamplesSubscriptions() {
	// Set current time function, so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(prettyJSONWriter{}, "ExampleEvents")

	handler := camunda.TaskHandlerFunc(func(ctx context.Context, completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
		l := logger.FromContext(ctx)
		l.Info().Msg("always logged")
	})

	h := middleware.SubscriptionsAddLogger(log)(handler)

	completeHandler := camunda.TaskCompleteFunc(func(ctx context.Context, taskID string) error {
		return nil
	})
	t := camunda.Task{}

	h.Handle(context.Background(), completeHandler, t)
}
