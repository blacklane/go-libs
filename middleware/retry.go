package middleware

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/blacklane/go-libs/x/events"
)

type RetriableError struct {
	Retriable bool
	Err       error
}

func (re *RetriableError) Error() string {
	if re.Err == nil {
		return "RetriableError.Err is nil"
	}
	return re.Err.Error()
}

func (re *RetriableError) Unwrap() error {
	return re.Err
}

func Retry(maxRetry int) events.Middleware {
	return func(handler events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			var err error

			for retries := 0; retries < maxRetry; retries++ {
				err = handler.Handle(ctx, e)
				if err == nil {
					return nil
				}

				var rerr *RetriableError
				if !errors.As(err, &rerr) {
					return err
				}
				if !rerr.Retriable {
					return fmt.Errorf("handler failed with non retriable error: %w", err)
				}

				// exponential backoff
				time.Sleep(time.Duration(math.Pow(2.0, float64(retries))*100) * time.Millisecond)
			}
			return fmt.Errorf("handler failed after %d retries: maximum retries exceeded: %w", maxRetry, err)
		})
	}
}
