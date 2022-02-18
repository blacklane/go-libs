# Middleware

TODO: write the readme

## Installation

```go
go get -u github.com/blacklane/go-libs/middleware
```

## Retry: Getting Started

`retry` is a middleware which for x/events. It will retry with exponential back-off algorithm until it exceeds max retry, which is passed by argument to `middleware.Retry`.

```go
maxRetryCount := 10

hb := events.HandlerBuilder{}
hb.AddHandler(handler)
hb.UseMiddleware(middleware.Retry(maxRetryCount))
```

If the handler returns a RetriableError, then the middleware will retry.

```go
type Handler struct {}

func (h Handler) Handle(ctx context.Context, e events.Event) error {
	err := processEvent(e)
	if err != nil {
		return RetriableError{
			Retriable: true,
			Err:       fmt.Errorf("fail to process event: %w", err),
		}
	}

	return nil
}
```

## Camunda Subscription Handler: Getting Started

The Camunda subscription handler is a middleware which allows the injection of the `context` from the service in ot the Camunda subscription handler.

```go
func CamundaSubscriptionsAddLogger(log logger.Logger) func(camunda.TaskHandler) camunda.TaskHandler {
	return func(next camunda.TaskHandler) camunda.TaskHandler {
		return camunda.TaskHandlerFunc(func(ctx context.Context, completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
			log = log.With().Logger()
			context := log.WithContext(ctx)
			next.Handle(context, completeFunc, t)
		})
	}
}
```
