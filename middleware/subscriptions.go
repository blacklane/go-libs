package middleware

import (
	"context"

	"github.com/blacklane/go-libs/camunda/v2"
	"github.com/blacklane/go-libs/logger"
)

func SubscriptionsAddLogger(log logger.Logger) func(camunda.TaskHandler) camunda.TaskHandler {
	return func(next camunda.TaskHandler) camunda.TaskHandler {
		return camunda.TaskHandlerFunc(func(ctx context.Context, completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
			log = log.With().Logger()
			context := log.WithContext(ctx)
			next.Handle(context, completeFunc, t)
		})
	}
}
