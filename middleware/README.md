# Middleware

TODO: write the readme

## Installation

```go
go get -u github.com/blacklane/go-libs/middleware
```

## Retry: getting Started

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
