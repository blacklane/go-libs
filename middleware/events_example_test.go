package middleware_test

import (
	"context"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"

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
		Key:     []byte("the_event_key"),
		TopicPartition: events.TopicPartition{
			Topic:     "the_topic",
			Partition: 1,
			Offset:    2,
		},
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
	//   "event_key": "the_event_key",
	//   "level": "info",
	//   "message": "always logged",
	//   "offset": 2,
	//   "partition": 1,
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "topic": "the_topic",
	//   "tracking_id": "the_tracking_id"
	// }
	// {
	//   "application": "ExampleEvents",
	//   "duration_ms": 0,
	//   "event": "event_to_be_logged",
	//   "event_key": "the_event_key",
	//   "level": "info",
	//   "message": "event_to_be_logged succeeded",
	//   "offset": 2,
	//   "partition": 1,
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "topic": "the_topic",
	//   "tracking_id": "the_tracking_id"
	// }
	// {
	//   "application": "ExampleEvents",
	//   "event": "event_to_not_log",
	//   "event_key": "",
	//   "level": "info",
	//   "message": "always logged",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "the_tracking_id",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "topic": "",
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
	//   "event_key": "",
	//   "level": "info",
	//   "message": "always logged",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "ExampleEvents_onlyRequestIDHeader",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "topic": "",
	//   "tracking_id": "ExampleEvents_onlyRequestIDHeader"
	// }
	// {
	//   "application": "ExampleEvents_onlyRequestIDHeader",
	//   "duration_ms": 0,
	//   "event": "event_to_be_logged",
	//   "event_key": "",
	//   "level": "info",
	//   "message": "event_to_be_logged succeeded",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "ExampleEvents_onlyRequestIDHeader",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "topic": "",
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
	//   "event_key": "",
	//   "level": "info",
	//   "message": "always logged",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "ExampleEvents_trackingIDAndRequestIDHeaders",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "topic": "",
	//   "tracking_id": "ExampleEvents_trackingIDAndRequestIDHeaders"
	// }
	// {
	//   "application": "ExampleEvents_trackingIDAndRequestIDHeaders",
	//   "duration_ms": 0,
	//   "event": "event_to_be_logged",
	//   "event_key": "",
	//   "level": "info",
	//   "message": "event_to_be_logged succeeded",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "ExampleEvents_trackingIDAndRequestIDHeaders",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "topic": "",
	//   "tracking_id": "ExampleEvents_trackingIDAndRequestIDHeaders"
	// }
}
