package middleware

import (
	"context"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	"github.com/blacklane/go-libs/otel"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/google/uuid"

	"github.com/blacklane/go-libs/middleware/internal/constants"
)

// Events wraps handler in middleware to provide:
// - a tracking ID in the context,
// - a logger in the context and log the "status" at the end of each handler,
// - a OTel span for handler.
func Events(handler events.Handler, log logger.Logger, eventName string) events.Handler {
	hb := events.HandlerBuilder{}
	hb.AddHandler(handler)
	hb.UseMiddleware(
		EventsTrackingID,
		logmiddleware.EventsAddLogger(log),
		otel.EventsAddOpenTelemetry(eventName),
		logmiddleware.EventsHandlerStatusLogger(eventName),
	)

	return hb.Build()[0]
}

// EventsTrackingID adds a tracking id to the event handler context.
// The tracking ID is, in order of preference,
// the value of the header constants.HeaderTrackingID OR
// the value of the header constants.HeaderRequestID OR
// a newly generated UUID.
func EventsTrackingID(next events.Handler) events.Handler {
	return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
		trackingID := eventsExtractTrackingID(e)

		ctx = tracking.SetContextID(ctx, trackingID)
		return next.Handle(ctx, e)
	})
}

func eventsExtractTrackingID(e events.Event) string {
	if trackingID, ok := e.Headers[constants.HeaderTrackingID]; ok {
		return trackingID
	}
	if requestID, ok := e.Headers[constants.HeaderRequestID]; ok {
		return requestID
	}

	// ignoring the error as the current implementation of uuid.NewUUID never
	// returns an error
	newUUID, _ := uuid.NewUUID()
	return newUUID.String()
}
