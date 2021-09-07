package middleware

import (
	"context"

	"github.com/blacklane/go-libs/x/events"
	"github.com/google/uuid"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/tracking/internal/constants"
)

// EventsAddTrackingID checks if an event has a tracking ID in the header
// if it does it sets it on the context, if it does not it generates a new one to set on the context
func EventsAddTrackingID(next events.Handler) events.Handler {
	return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
		trackingID := eventsExtractTrackingID(e)

		ctx = tracking.SetContextID(ctx, trackingID)
		return next.Handle(ctx, e)
	})
}

func eventsExtractTrackingID(e events.Event) string {
	trackingID := e.Headers[constants.HeaderTrackingID]
	if trackingID == "" {
		trackingID = e.Headers[constants.HeaderRequestID]
		if trackingID == "" {
			uuid, err := uuid.NewUUID()
			if err != nil {
				return ""
			}
			trackingID = uuid.String()
		}
	}

	return trackingID
}
