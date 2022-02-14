package middleware_test

import (
	"context"

	"github.com/blacklane/go-libs/camunda/v2"
	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/middleware"
)

func ExamplesSubscriptions() {
	log := logger.New(prettyJSONWriter{}, "ExampleEvents")

	handler := camunda.TaskHandlerFunc(func(ctx context.Context, completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
		l := logger.FromContext(ctx)
		l.Info().Msg("always logged")
	})

	h := middleware.CamundaSubscriptionsAddLogger(log)(handler)

	completeHandler := camunda.TaskCompleteFunc(func(ctx context.Context, taskID string) error {
		return nil
	})
	t := camunda.Task{}

	h.Handle(context.Background(), completeHandler, t)
}
