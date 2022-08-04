package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
)

func ExampleEventsAddLogger() {
	// Set current time function so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(prettyJSONWriter{}, "ExampleEventsAddLogger")

	h := events.HandlerFunc(func(ctx context.Context, _ events.Event) error {
		l := logger.FromContext(ctx)
		l.Info().Msg("Hello, Gophers from events")
		return nil
	})

	m := EventsAddLogger(log)
	hh := m(h)

	_ = hh.Handle(context.Background(), events.Event{})

	// Output:
	// {
	//   "application": "ExampleEventsAddLogger",
	//   "level": "info",
	//   "message": "Hello, Gophers from events",
	//   "timestamp": "2009-11-10T23:00:00.000Z"
	// }
}

func ExampleEventsHandlerStatusLogger_success() {
	// Set current time function so we can control the logged timestamp and duration
	timeNowCalled := false
	logger.SetNowFunc(func() time.Time {
		now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		if timeNowCalled {
			now = now.Add(time.Second)
		}
		timeNowCalled = true
		return now
	})

	ctx := tracking.SetContextID(context.Background(), "tracking_id-ExampleEventsLogger_Success")

	log := logger.New(prettyJSONWriter{}, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsAddLogger(log), EventsHandlerStatusLogger())
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return nil }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{
		Key:     []byte("event_key_here"),
		Payload: []byte(`{"event":"event_name_here"}`),
		TopicPartition: events.TopicPartition{
			Topic:     "topic_here",
			Partition: 1,
			Offset:    2,
		},
	})

	// Output:
	// {
	//   "application": "ExampleEventsLogger",
	//   "duration_ms": 1000,
	//   "event": "event_name_here",
	//   "event_key": "event_key_here",
	//   "level": "info",
	//   "message": "event_name_here succeeded",
	//   "offset": 2,
	//   "partition": 1,
	//   "request_id": "tracking_id-ExampleEventsLogger_Success",
	//   "timestamp": "2009-11-10T23:00:01.000Z",
	//   "topic": "topic_here",
	//   "tracking_id": "tracking_id-ExampleEventsLogger_Success"
	// }
}

func ExampleEventsHandlerStatusLogger_onlyLogCertainEvents() {
	// Set current time function so we can control the logged timestamp and duration
	timeNowCalled := false
	logger.SetNowFunc(func() time.Time {
		now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		if timeNowCalled {
			now = now.Add(time.Second)
		}
		timeNowCalled = true
		return now
	})

	ctx := tracking.SetContextID(context.Background(), "tracking_id-ExampleEventsLogger_Success")

	log := logger.New(prettyJSONWriter{}, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsAddLogger(log), EventsHandlerStatusLogger("log_event"))
	hb.AddHandler(
		events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			log := logger.FromContext(ctx)
			log.Info().Msgf("Log from handler")
			return nil
		}))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"log_event"}`)})
	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"do_not_log_event"}`)})

	// Output:
	// {
	//   "application": "ExampleEventsLogger",
	//   "event": "log_event",
	//   "event_key": "",
	//   "level": "info",
	//   "message": "Log from handler",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "tracking_id-ExampleEventsLogger_Success",
	//   "timestamp": "2009-11-10T23:00:01.000Z",
	//   "topic": "",
	//   "tracking_id": "tracking_id-ExampleEventsLogger_Success"
	// }
	// {
	//   "application": "ExampleEventsLogger",
	//   "duration_ms": 1000,
	//   "event": "log_event",
	//   "event_key": "",
	//   "level": "info",
	//   "message": "log_event succeeded",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "tracking_id-ExampleEventsLogger_Success",
	//   "timestamp": "2009-11-10T23:00:01.000Z",
	//   "topic": "",
	//   "tracking_id": "tracking_id-ExampleEventsLogger_Success"
	// }
	// {
	//   "application": "ExampleEventsLogger",
	//   "event": "do_not_log_event",
	//   "event_key": "",
	//   "level": "info",
	//   "message": "Log from handler",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "tracking_id-ExampleEventsLogger_Success",
	//   "timestamp": "2009-11-10T23:00:01.000Z",
	//   "topic": "",
	//   "tracking_id": "tracking_id-ExampleEventsLogger_Success"
	// }
}
func ExampleEventsHandlerStatusLogger_failure() {
	// Set current time function so we can control the logged timestamp and duration
	timeNowCalled := false
	logger.SetNowFunc(func() time.Time {
		now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		if timeNowCalled {
			now = now.Add(time.Second)
		}
		timeNowCalled = true
		return now
	})

	ctx := tracking.SetContextID(context.Background(), "tracking_id-ExampleEventsLogger_Failure")

	log := logger.New(prettyJSONWriter{}, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsAddLogger(log), EventsHandlerStatusLogger())
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return errors.New("bad") }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"event_name_here"}`)})

	// Output:
	// {
	//   "application": "ExampleEventsLogger",
	//   "duration_ms": 1000,
	//   "error": "bad",
	//   "event": "event_name_here",
	//   "event_key": "",
	//   "level": "error",
	//   "message": "event_name_here failed",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "tracking_id-ExampleEventsLogger_Failure",
	//   "timestamp": "2009-11-10T23:00:01.000Z",
	//   "topic": "",
	//   "tracking_id": "tracking_id-ExampleEventsLogger_Failure"
	// }
}

func ExampleEventsHandlerStatusLoggerWithNameFn() {
	// Set current time function so we can control the logged timestamp and duration
	timeNowCalled := false
	logger.SetNowFunc(func() time.Time {
		now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		if timeNowCalled {
			now = now.Add(time.Second)
		}
		timeNowCalled = true
		return now
	})

	ctx := tracking.SetContextID(context.Background(), "tracking_id-ExampleEventsHandlerStatusLoggerWithNameFn")

	log := logger.New(prettyJSONWriter{}, "ExampleEventsHandlerStatusLoggerWithNameFn")

	hb := events.HandlerBuilder{}
	nameFn := func(e events.Event) string {
		type ePayload struct {
			Name string `json:"name"`
		}
		payload := &ePayload{}

		if err := json.Unmarshal(e.Payload, payload); err != nil {
			return "name_error"
		}

		return payload.Name
	}
	hb.UseMiddleware(EventsAddLogger(log), EventsHandlerStatusLoggerWithNameFn(nameFn))
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return nil }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"name":"event_name_here"}`)})

	// Output:
	// {
	//   "application": "ExampleEventsHandlerStatusLoggerWithNameFn",
	//   "duration_ms": 1000,
	//   "event": "event_name_here",
	//   "event_key": "",
	//   "level": "info",
	//   "message": "event_name_here succeeded",
	//   "offset": 0,
	//   "partition": 0,
	//   "request_id": "tracking_id-ExampleEventsHandlerStatusLoggerWithNameFn",
	//   "timestamp": "2009-11-10T23:00:01.000Z",
	//   "topic": "",
	//   "tracking_id": "tracking_id-ExampleEventsHandlerStatusLoggerWithNameFn"
	// }
}
