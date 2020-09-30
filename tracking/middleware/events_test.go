package middleware

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/tracking/internal/constants"
	"github.com/blacklane/go-libs/x/events"
)

func TestEventsAddTrackingIDCreatesIDWhenEventHeaderEmpty(t *testing.T) {
	e := events.Event{
		Headers: events.Header(map[string]string{constants.HeaderRequestID: ""}),
	}

	testHandler := EventsAddTrackingID(events.Handler(
		events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			if tracking.IDFromContext(ctx) == "" {
				t.Errorf("request id in context empty")
			}
			return nil
		})))

	err := testHandler.Handle(context.Background(), e)
	if err != nil {
		t.Errorf("could not successfully handle: %v", err)
	}
}

func TestEventsAddTrackingIDDoesNotChangeTrackingIDIfAlreadyPresent(t *testing.T) {
	testId := "goodid"
	for _, fieldName := range []string{constants.HeaderRequestID, constants.HeaderTrackingID} {
		e := events.Event{
			Headers: events.Header(map[string]string{fieldName: testId}),
		}

		testHandler :=  EventsAddTrackingID(events.Handler(
			events.HandlerFunc(func(ctx context.Context, e events.Event) error {
				got := tracking.IDFromContext(ctx)
				if !cmp.Equal(got, testId) {
					t.Errorf("field name: %s, want: %v, got: %v", fieldName, testId, got)
				}
				return nil
			})))

		err := testHandler.Handle(context.Background(), e)
		if err != nil {
			t.Errorf("field name: %s, could not successfully handle: %v", fieldName, err)
		}
	}
}
