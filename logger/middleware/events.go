package middleware

import (
	"context"
	"encoding/json"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/rs/zerolog"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
)

// EventsAddLogger adds the logger into the context.
func EventsAddLogger(log logger.Logger) events.Middleware {
	return func(next events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			return next.Handle(log.WithContext(ctx), e)
		})
	}
}

// EventsHandlerStatusLogger produces a log line after every event handled with
// the status (success or failure), tracking id, duration, event name.
// If it failed the error is added to the log line as well.
// It also updates the logger in the context adding the tracking id present in
// the context.
//
// To extract the event name it considers the event.Payload is a JSON with a top
// level key 'event' and uses its value as the event name.
// It is possible to pass a custom function to extract the event name,
// check EventsLoggerWithNameFn.
func EventsHandlerStatusLogger() events.Middleware {
	return EventsHandlerStatusLoggerWithNameFn(eventName)
}

// EventsLoggerWithNameFn is the same as EventsHandlerStatusLogger, but using a custom
// function to extract the event name.
func EventsHandlerStatusLoggerWithNameFn(eventNameFn func(e events.Event) string) events.Middleware {
	return func(next events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) (err error) {
			startTime := logger.Now()

			if eventNameFn == nil {
				eventNameFn = eventName
			}
			evName := eventNameFn(e)

			l := logger.FromContext(ctx)
			trackingID := tracking.IDFromContext(ctx)
			logFields := map[string]interface{}{
				internal.FieldTrackingID: trackingID,
				internal.FieldRequestID:  trackingID,
				internal.FieldEvent:      evName,
			}
			l.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Fields(logFields)
			})

			defer func() {
				status := "succeeded"

				duration := logger.Now().Sub(startTime)

				zlEvent := l.Info().
					Dur(internal.FieldDuration, duration)

				if err != nil {
					status = "failed"
					zlEvent.Err(err)
				}
				zlEvent.Msgf("%s %s", evName, status)
			}()
			return next.Handle(ctx, e)
		})
	}
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
