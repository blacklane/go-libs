package middleware_test

import (
	"context"
	"testing"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/google/go-cmp/cmp"

	"github.com/blacklane/go-libs/middleware"
	"github.com/blacklane/go-libs/middleware/internal/constants"
)

func ExampleEvents() {
	trackingID := "the_tracking_id"

	// Set current time function, so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(prettyJSONWriter{}, "ExampleEvents")

	h := events.HandlerFunc(func(ctx context.Context, _ events.Event) error {
		l := logger.FromContext(ctx)
		l.Info().Msg("always logged")
		return nil
	})

	eventNameToLog := "event_to_be_logged"
	hh := middleware.Events(h, log, eventNameToLog)

	eventToLog := events.Event{
		Headers: map[string]string{constants.HeaderTrackingID: trackingID},
		Payload: []byte(`{"event":"` + eventNameToLog + `"}`),
	}
	eventToNotLog := events.Event{
		Headers: map[string]string{constants.HeaderTrackingID: trackingID},
		Payload: []byte(`{"event":"event_to_not_log"}`),
	}

	_ = hh.Handle(context.Background(), eventToLog)
	_ = hh.Handle(context.Background(), eventToNotLog)

	// Output:
	// {
	//   "application": "ExampleEvents",
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "always logged",
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "the_tracking_id"
	// }
	// {
	//   "application": "ExampleEvents",
	//   "duration_ms": 0,
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "event_to_be_logged succeeded",
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "the_tracking_id"
	// }
	// {
	//   "application": "ExampleEvents",
	//   "event": "event_to_not_log",
	//   "level": "info",
	//   "message": "always logged",
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "the_tracking_id"
	// }
}

func ExampleEvents_onlyRequestIDHeader() {
	trackingID := "ExampleEvents_onlyRequestIDHeader"

	// Set current time function, so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(prettyJSONWriter{}, "ExampleEvents_onlyRequestIDHeader")

	h := events.HandlerFunc(func(ctx context.Context, _ events.Event) error {
		l := logger.FromContext(ctx)
		l.Info().Msg("always logged")
		return nil
	})

	eventNameToLog := "event_to_be_logged"
	hh := middleware.Events(h, log, eventNameToLog)

	eventToLog := events.Event{
		Headers: map[string]string{constants.HeaderRequestID: trackingID},
		Payload: []byte(`{"event":"` + eventNameToLog + `"}`),
	}

	_ = hh.Handle(context.Background(), eventToLog)

	// Output:
	// {
	//   "application": "ExampleEvents_onlyRequestIDHeader",
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "always logged",
	//   "request_id": "ExampleEvents_onlyRequestIDHeader",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleEvents_onlyRequestIDHeader"
	// }
	// {
	//   "application": "ExampleEvents_onlyRequestIDHeader",
	//   "duration_ms": 0,
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "event_to_be_logged succeeded",
	//   "request_id": "ExampleEvents_onlyRequestIDHeader",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleEvents_onlyRequestIDHeader"
	// }
}

func ExampleEvents_trackingIDAndRequestIDHeaders() {
	trackingID := "ExampleEvents_trackingIDAndRequestIDHeaders"

	// Set current time function, so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(prettyJSONWriter{}, "ExampleEvents_trackingIDAndRequestIDHeaders")

	h := events.HandlerFunc(func(ctx context.Context, _ events.Event) error {
		l := logger.FromContext(ctx)
		l.Info().Msg("always logged")
		return nil
	})

	eventNameToLog := "event_to_be_logged"
	hh := middleware.Events(h, log, eventNameToLog)

	eventToLog := events.Event{
		Headers: map[string]string{constants.HeaderRequestID: trackingID},
		Payload: []byte(`{"event":"` + eventNameToLog + `"}`),
	}

	_ = hh.Handle(context.Background(), eventToLog)

	// Output:
	// {
	//   "application": "ExampleEvents_trackingIDAndRequestIDHeaders",
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "always logged",
	//   "request_id": "ExampleEvents_trackingIDAndRequestIDHeaders",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleEvents_trackingIDAndRequestIDHeaders"
	// }
	// {
	//   "application": "ExampleEvents_trackingIDAndRequestIDHeaders",
	//   "duration_ms": 0,
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "event_to_be_logged succeeded",
	//   "request_id": "ExampleEvents_trackingIDAndRequestIDHeaders",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleEvents_trackingIDAndRequestIDHeaders"
	// }
}

func TestEventsTrackingIDCreatesIDWhenEventHeaderEmpty(t *testing.T) {
	e := events.Event{
		Headers: events.Header(map[string]string{constants.HeaderRequestID: ""}),
	}

	testHandler := middleware.EventsTrackingID(events.Handler(
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
