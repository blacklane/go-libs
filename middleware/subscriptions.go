package middleware

import (
	"context"

	"github.com/blacklane/go-libs/camunda"
	"github.com/blacklane/go-libs/logger"
)

func SubscriptionsAddLogger(log logger.Logger) func(camunda.TaskHandler) camunda.TaskHandler {
	return func(next camunda.TaskHandler) camunda.TaskHandler {
		h := camunda.TaskHandlerFunc(func(ctx context.Context, completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
			log = log.With().Logger() // Creates a sub logger so all requests won't share the same logger instance
			context := log.WithContext(ctx)
			next.Handle(context, completeFunc, t)
		})
		return h
	}
}
