package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/blacklane/go-libs/x/events"
)

type mockHandler struct {
	Called      int
	handlerFunc events.HandlerFunc
}

func (m *mockHandler) Handle(ctx context.Context, e events.Event) error {
	m.Called++
	return m.handlerFunc.Handle(ctx, e)
}

func TestRetryMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		wantHandlerCalled int
		handleReturnErr   error
		wantErr           error
	}{
		{"RetriableError", 3, &RetriableError{Retriable: true, Err: errors.New("SampleError")}, errors.New("handler failed after 3 retries: maximum retries exceeded: SampleError")},
		{"NonRetriableError", 1, &RetriableError{Retriable: false, Err: errors.New("SampleError")}, errors.New("handler failed with non retriable error: SampleError")},
		{"NormalError", 1, errors.New("SampleError"), errors.New("SampleError")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantHandlerCalled := tt.wantHandlerCalled
			maxRetries := 3

			handler := &mockHandler{
				handlerFunc: events.HandlerFunc(func(ctx context.Context, e events.Event) error {
					return tt.handleReturnErr
				}),
			}

			err := Retry(maxRetries)(handler).Handle(context.Background(), events.Event{})

			if handler.Called != wantHandlerCalled {
				t.Errorf("Handle() is expected to be called %d times but it called %d times", wantHandlerCalled, handler.Called)
			}

			if err.Error() != tt.wantErr.Error() {
				t.Errorf("Handle() is expected to return err: %v, but it returns err: %v", tt.wantErr, err)
			}
		})
	}
}

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		name           string
		retryCount     int
		wantMillsecond time.Duration
	}{
		{"retry 0", 0, 100 * time.Millisecond},
		{"retry 1", 1, 200 * time.Millisecond},
		{"retry 2", 2, 400 * time.Millisecond},
		{"retry 3", 3, 800 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exponentialBackoff(tt.retryCount)

			if tt.wantMillsecond != got {
				t.Errorf("exponentialBackoff is expected to return: %d, but it returned: %d", tt.wantMillsecond, got)
			}
		},
		)
	}
}
