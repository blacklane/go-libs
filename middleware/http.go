package middleware

import (
	"net/http"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	"github.com/blacklane/go-libs/otel"
	"github.com/blacklane/go-libs/tracking"
	"github.com/google/uuid"

	"github.com/blacklane/go-libs/middleware/internal/constants"
)

// HTTP returns the composition of all Blacklane's middleware.
func HTTP(serviceName, handlerName, path string, log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := otel.HTTPMiddleware(serviceName, handlerName, path)(next)
			h = logmiddleware.HTTPRequestLogger([]string{})(h)
			h = logmiddleware.HTTPAddLogger(log)(h)
			h = HTTPTrackingID(h)

			h.ServeHTTP(w, r)
		})
	}
}

// HTTPWithBodyFilter returns the composition of all Blacklane's middleware plus body.
func HTTPWithBodyFilter(serviceName, handlerName, path string, filterKeys []string, log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := otel.HTTPMiddleware(serviceName, handlerName, path)(next)
			h = logmiddleware.HTTPRequestLogger([]string{})(h)
			h = logmiddleware.HTTPAddLogger(log)(h)
			h = logmiddleware.HTTPAddBodyFilters(filterKeys)(h)
			h = HTTPTrackingID(h)

			h.ServeHTTP(w, r)
		})
	}
}

// HTTPFunc is a helper function to use ordinary functions instead of http.Handler as HTTP requires.
// It makes the necessary type conversions and delegates to HTTP.
// Check HTTP for details.
// TODO: add tests.
func HTTPFunc(serviceName, handlerName, path string, log logger.Logger) func(func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return HTTP(serviceName, handlerName, path, log)(http.HandlerFunc(next)).ServeHTTP
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
	if trackingID := r.Header.Get(constants.HeaderTrackingID); trackingID != "" {
		return trackingID
	}
	if requestID := r.Header.Get(constants.HeaderRequestID); requestID != "" {
		return requestID
	}

	// ignoring the error as the current implementation of uuid.NewUUID never
	// returns an error and there isn't much we can do if it'd fail anyway.
	newUUID, _ := uuid.NewUUID()
	return newUUID.String()
}
