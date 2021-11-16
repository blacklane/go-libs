package middleware_test

import (
	"context"
	"testing"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/google/go-cmp/cmp"

	"github.com/blacklane/go-libs/middleware"
	"github.com/blacklane/go-libs/middleware/internal/constants"
)

func TestEventsTrackingIDCreatesIDWhenEventHeaderEmpty(t *testing.T) {
	e := events.Event{
		Headers: events.Header(
			map[string]string{constants.HeaderRequestID: ""}),
	}

	testHandler := middleware.EventsTrackingID(events.Handler(
		events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			if tracking.IDFromContext(ctx) == "" {
				t.Errorf("tracking id in context is empty")
			}
			return nil
		})))

	err := testHandler.Handle(context.Background(), e)
	if err != nil {
		t.Errorf("Handle returned an error: %v", err)
	}
}

func TestEventsTrackingIDDoesNotChangeTrackingIDIfAlreadyPresent(t *testing.T) {
	testId := "goodid"
	for _, headerName := range []string{constants.HeaderRequestID, constants.HeaderTrackingID} {
		e := events.Event{
			Headers: events.Header(map[string]string{headerName: testId}),
		}

		testHandler := middleware.EventsTrackingID(events.Handler(
			events.HandlerFunc(func(ctx context.Context, e events.Event) error {
				got := tracking.IDFromContext(ctx)
				if !cmp.Equal(got, testId) {
					t.Errorf("field name: %s, want: %v, got: %v", headerName, testId, got)
				}
				return nil
			})))

		err := testHandler.Handle(context.Background(), e)
		if err != nil {
			t.Errorf("field name: %s, could not successfully handle: %v", headerName, err)
		}
	}
}
