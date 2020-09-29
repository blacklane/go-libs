package middleware

import (
	"context"

	"github.com/google/uuid"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/tracking/internal/constants"
	"github.com/blacklane/go-libs/x/events"
)

// EventsAddTrackingID checks if an event has a tracking ID in the header
// if it does it sets it on the context, if it does not it generates a new one to set on the context
func EventsAddTrackingID(next events.Handler) events.Handler {
		return events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			trackingID := e.Headers[constants.HeaderTrackingID]
			if trackingID == "" {
				trackingID = e.Headers[constants.HeaderRequestID]
			}
			if trackingID == "" {
				uuid, err := uuid.NewUUID()
				if err != nil {
					logger.FromContext(ctx).Err(err).Msg("could not generate uuid not setting trackingID in context")
					return next.Handle(ctx, e)
				}
				trackingID = uuid.String()
			}
			ctx = tracking.SetContextID(ctx, trackingID)
			return next.Handle(ctx, e)
		})
	}
