package tracing

import (
	"context"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"

	"github.com/blacklane/go-libs/tracing/internal/constants"
)

func ExampleEventsAddDefault() {
	trackingID := "the_tracking_id"

	// Set current time function so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(prettyJSONWriter{}, "ExampleEventsAddAll")

	h := events.HandlerFunc(func(ctx context.Context, _ events.Event) error {
		l := logger.FromContext(ctx)
		l.Info().Msg("always logged")
		return nil
	})

	eventNameToLog := "event_to_be_logged"
	hh := EventsAddDefault(h, log, eventNameToLog)

	eventToLog := events.Event{
		Headers: map[string]string{"X-Tracking-Id": trackingID},
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
	//   "application": "ExampleEventsAddAll",
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "always logged",
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "tracking_id": "the_tracking_id"
	// }
	// {
	//   "application": "ExampleEventsAddAll",
	//   "duration_ms": 0,
	//   "event": "event_to_be_logged",
	//   "level": "info",
	//   "message": "event_to_be_logged succeeded",
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "tracking_id": "the_tracking_id"
	// }
	// {
	//   "application": "ExampleEventsAddAll",
	//   "event": "event_to_not_log",
	//   "level": "info",
	//   "message": "always logged",
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "tracking_id": "the_tracking_id"
	// }
}
