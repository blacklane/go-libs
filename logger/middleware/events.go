package middleware

import (
	"context"
	"encoding/json"

	"github.com/blacklane/go-libs/tracking"
	trackingMidleware "github.com/blacklane/go-libs/tracking/middleware"
	"github.com/blacklane/go-libs/x/events"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
)

// EventsAddAll adds the necessary middleware for:
//   - have tracking id in the context (read from the headers or a new one),
//   - have a logger.Logger with tracking id and all required fields in the context,
//   - log at the end of handler if it succeeded or failed and how log it took.
// For more details, check the middleware used:
// - github.com/blacklane/go-libs/tracking/middleware.EventsAddTrackingID
// - middleware.EventsAddLogger
// - middleware.EventsHandlerStatusLogger
func EventsAddAll(handler events.Handler, log logger.Logger, eventNames ...string) events.Handler {
	hb := events.HandlerBuilder{}
	hb.AddHandler(handler)
	hb.UseMiddleware(
		trackingMidleware.EventsAddTrackingID,
		EventsAddLogger(log),
		EventsHandlerStatusLogger(eventNames...))

	return hb.Build()[0]
}

// EventsAddLogger adds the logger into the context.
func EventsAddLogger(log logger.Logger) events.Middleware {
	return func(next events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			log := log.With().Logger()
			ctx = log.WithContext(ctx)
			return next.Handle(log.WithContext(ctx), e)
		})
	}
}

// EventsHandlerStatusLogger produces a log line after every event handled with
// the status (success or failure), tracking id, duration, event name.
// If it failed the error is added to the log line as well.
// It also updates the logger in the context adding the tracking id present in
// the context.
// If one or more eventNames are provided, only events matching these names are
// logged. If none os provided, all events are logged.
//
// To extract the event name it considers the event.Payload is a JSON with a top
// level key 'event' and uses its value as the event name.
// It is possible to pass a custom function to extract the event name,
// check EventsLoggerWithNameFn.
func EventsHandlerStatusLogger(eventNames ...string) events.Middleware {
	return EventsHandlerStatusLoggerWithNameFn(eventName, eventNames...)
}

// EventsLoggerWithNameFn is the same as EventsHandlerStatusLogger, but using a custom
// function to extract the event name.
func EventsHandlerStatusLoggerWithNameFn(
	eventNameFn func(e events.Event) string,
	eventNames ...string) events.Middleware {
	return func(next events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) (err error) {
			startTime := logger.Now()

			if eventNameFn == nil {
				eventNameFn = eventName
			}
			evName := eventNameFn(e)

			log := *logger.FromContext(ctx)
			trackingID := tracking.IDFromContext(ctx)
			logFields := map[string]interface{}{
				internal.FieldTrackingID: trackingID,
				internal.FieldRequestID:  trackingID,
				internal.FieldEvent:      evName,
			}

			log = log.With().Fields(logFields).Logger()
			ctx = log.WithContext(ctx)

			if !logEvent(evName, eventNames...) {
				return next.Handle(ctx, e)
			}

			defer func() {
				zlEvent := log.Info()
				status := "succeeded"

				duration := logger.Now().Sub(startTime)

				if err != nil {
					zlEvent = log.Error()
					status = "failed"
					zlEvent.Err(err)
				}
				zlEvent.
					Dur(internal.FieldDuration, duration).
					Msgf("%s %s", evName, status)
			}()
			return next.Handle(ctx, e)
		})
	}
}

func logEvent(event string, eventNamesToLog ...string) bool {
	// We assume all events should be logged unless otherwise specified.
	if len(eventNamesToLog) == 0 {
		return true
	}

	for _, e := range eventNamesToLog {
		if event == e {
			return true
		}
	}

	return false
}

func eventName(e events.Event) string {
	type ePayload struct {
		Name string `json:"event"`
	}
	payload := &ePayload{}

	if err := json.Unmarshal(e.Payload, payload); err != nil {
		return "name_error"
	}

	return payload.Name
}
