package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"os"
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

	log := logger.New(os.Stdout, "ExampleEventsAddLogger")

	h := events.HandlerFunc(func(ctx context.Context, _ events.Event) error {
		l := logger.FromContext(ctx)
		l.Info().Msg("Hello, Gophers from events")
		return nil
	})

	m := EventsAddLogger(log)
	hh := m(h)

	_ = hh.Handle(context.Background(), events.Event{})

	// Output:
	// {"level":"info","application":"ExampleEventsAddLogger","timestamp":"2009-11-10T23:00:00Z","message":"Hello, Gophers from events"}
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

	log := logger.New(os.Stdout, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsHandlerStatusLogger(), EventsAddLogger(log))
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return nil }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"event_name_here"}`)})

	// Output:
	// {"level":"info","application":"ExampleEventsLogger","event":"event_name_here","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","duration_ms":1000,"timestamp":"2009-11-10T23:00:01Z","message":"event_name_here succeeded"}
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

	log := logger.New(os.Stdout, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsHandlerStatusLogger("log_event"), EventsAddLogger(log))
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return nil }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"log_event"}`)})
	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"do_not_log_event"}`)})

	// Output:
	// {"level":"info","application":"ExampleEventsLogger","event":"log_event","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","duration_ms":1000,"timestamp":"2009-11-10T23:00:01Z","message":"log_event succeeded"}
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

	log := logger.New(os.Stdout, "ExampleEventsLogger")

	hb := events.HandlerBuilder{}
	hb.UseMiddleware(EventsHandlerStatusLogger(), EventsAddLogger(log))
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return errors.New("bad") }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"event_name_here"}`)})

	// Output:
	// {"level":"error","application":"ExampleEventsLogger","event":"event_name_here","request_id":"tracking_id-ExampleEventsLogger_Failure","tracking_id":"tracking_id-ExampleEventsLogger_Failure","error":"bad","duration_ms":1000,"timestamp":"2009-11-10T23:00:01Z","message":"event_name_here failed"}
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

	ctx := tracking.SetContextID(context.Background(), "tracking_id-ExampleEventsLogger_Success")

	log := logger.New(os.Stdout, "ExampleEventsLogger")

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
	hb.UseMiddleware(EventsHandlerStatusLoggerWithNameFn(nameFn), EventsAddLogger(log))
	hb.AddHandler(
		events.HandlerFunc(func(context.Context, events.Event) error { return nil }))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"name":"event_name_here"}`)})

	// Output:
	// {"level":"info","application":"ExampleEventsLogger","event":"event_name_here","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","duration_ms":1000,"timestamp":"2009-11-10T23:00:01Z","message":"event_name_here succeeded"}
}

func ExampleLogWithRequiredFieldsForAllEvents() {
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
			log.Info().Msgf("Handler function execution with tracking id = %v", tracking.IDFromContext(ctx))
			return nil
		}))

	h := hb.Build()[0]

	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"log_event"}`)})
	_ = h.Handle(ctx, events.Event{Payload: []byte(`{"event":"do_not_log_event"}`)})

	// Output:
	// {"level":"info","application":"ExampleEventsLogger","event":"log_event","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","timestamp":"2009-11-10T23:00:01Z","message":"Handler function execution with tracking id = tracking_id-ExampleEventsLogger_Success"}
	// {"level":"info","application":"ExampleEventsLogger","event":"log_event","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","duration_ms":1000,"timestamp":"2009-11-10T23:00:01Z","message":"log_event succeeded"}
	// {"level":"info","application":"ExampleEventsLogger","event":"do_not_log_event","request_id":"tracking_id-ExampleEventsLogger_Success","tracking_id":"tracking_id-ExampleEventsLogger_Success","timestamp":"2009-11-10T23:00:01Z","message":"Handler function execution with tracking id = tracking_id-ExampleEventsLogger_Success"}
}
