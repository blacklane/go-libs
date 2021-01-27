package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"

	"github.com/blacklane/go-libs/logger"
)

func ExampleEventsAddLogger() {
	// Set current time function so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(prettyWriter, "ExampleEventsAddLogger")

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
	//   "timestamp": "2009-11-10T23:00:00Z"
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

	log := logger.New(prettyWriter, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsAddLogger(log), EventsHandlerStatusLogger())
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return nil }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"event_name_here"}`)})

	// Output:
	// {
	//   "application": "ExampleEventsLogger",
	//   "duration_ms": 1000,
	//   "event": "event_name_here",
	//   "level": "info",
	//   "message": "event_name_here succeeded",
	//   "request_id": "tracking_id-ExampleEventsLogger_Success",
	//   "timestamp": "2009-11-10T23:00:01Z",
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

	log := logger.New(prettyWriter, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsAddLogger(log), EventsHandlerStatusLogger("log_event"))
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return nil }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"log_event"}`)})
	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"do_not_log_event"}`)})

	// Output:
	// {
	//   "application": "ExampleEventsLogger",
	//   "duration_ms": 1000,
	//   "event": "log_event",
	//   "level": "info",
	//   "message": "log_event succeeded",
	//   "request_id": "tracking_id-ExampleEventsLogger_Success",
	//   "timestamp": "2009-11-10T23:00:01Z",
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

	log := logger.New(prettyWriter, "ExampleEventsLogger")

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
	//   "level": "error",
	//   "message": "event_name_here failed",
	//   "request_id": "tracking_id-ExampleEventsLogger_Failure",
	//   "timestamp": "2009-11-10T23:00:01Z",
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

	log := logger.New(prettyWriter, "ExampleEventsHandlerStatusLoggerWithNameFn")

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
	//   "level": "info",
	//   "message": "event_name_here succeeded",
	//   "request_id": "tracking_id-ExampleEventsHandlerStatusLoggerWithNameFn",
	//   "timestamp": "2009-11-10T23:00:01Z",
	//   "tracking_id": "tracking_id-ExampleEventsHandlerStatusLoggerWithNameFn"
	// }
}

func ExampleEventsHandlerStatusLogger_loggerFieldsSetForAllEvents() {
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

	log := logger.New(os.Stdout, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsHandlerStatusLogger("log_event"), EventsAddLogger(log))
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
	// {"level":"info","application":"ExampleEventsLogger","event":"log_event","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","timestamp":"2009-11-10T23:00:01.000Z","message":"Log from handler"}
	// {"level":"info","application":"ExampleEventsLogger","event":"log_event","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","duration_ms":1000,"timestamp":"2009-11-10T23:00:01.000Z","message":"log_event succeeded"}
	// {"level":"info","application":"ExampleEventsLogger","event":"do_not_log_event","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","timestamp":"2009-11-10T23:00:01.000Z","message":"Log from handler"}
}
