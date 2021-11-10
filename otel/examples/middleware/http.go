package middleware

import (
	"net/http"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	"github.com/blacklane/go-libs/otel"
	"github.com/blacklane/go-libs/tracking"
	"github.com/google/uuid"
)

// HTTP returns the composition of HTTPGenericMiddleware and
// HTTPMiddleware. Therefore, it should be applied to each http.Handler
// individually.
func HTTP(serviceName, handlerName, path string, log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := otel.HTTPMiddleware(serviceName, handlerName, path)(next)
			h = logmiddleware.HTTPRequestLogger(skipRoutes)(h)
			h = logmiddleware.HTTPAddLogger(log)(h)
			h = HTTPTrackingID(h)

			h.ServeHTTP(w, r)
		})
	}
}

// HTTPTrackingID adds a tracking id to the request context.
// The tracking ID is, in order of preference,
// the value of the header constants.HeaderTrackingID OR
// the value of the header constants.HeaderRequestID OR
// a newly generated UUID.
func HTTPTrackingID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackingID := newTrackingID(r)

		ctx := tracking.SetContextID(r.Context(), trackingID)
		rr := r.WithContext(ctx)

		next.ServeHTTP(w, rr)
	})
}

func newTrackingID(r *http.Request) string {
	if trackingID := r.Header.Get(HeaderTrackingID); trackingID != "" {
		return trackingID
	}
	if requestID := r.Header.Get(HeaderRequestID); requestID != "" {
		return requestID
	}

	// ignoring the error as the current implementation of uuid.NewUUID never
	// returns an error and there isn't much we can do if it'd fail anyway.
	newUUID, _ := uuid.NewUUID()
	return newUUID.String()
}
