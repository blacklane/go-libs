package uhttp

import (
	"net/http"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	"github.com/blacklane/go-libs/otel"
	"github.com/blacklane/go-libs/tracking"
	trackmiddleware "github.com/blacklane/go-libs/tracking/middleware"
	"github.com/google/uuid"

	"github.com/blacklane/go-libs/uhttp/internal/constants"
)

// HTTPAllMiddleware returns the composition of HTTPGenericMiddleware and
// HTTPMiddleware. Therefore it should be applied to each http.Handler
// individually.
func HTTPAllMiddleware(serviceName, handlerName, path string, log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := otel.HTTPMiddleware(serviceName, handlerName, path)(next)
			h = HTTPGenericMiddleware(log, skipRoutes...)(h)

			h.ServeHTTP(w, r)
		})
	}
}

// HTTPGenericMiddleware adds all middlewares which aren't specific to a handler.
// They are:
// - github.com/blacklane/go-libs/tracking/middleware.TrackingID
// - github.com/blacklane/go-libs/logger/middleware.HTTPAddLogger
// - github.com/blacklane/go-libs/logger/middleware.HTTPRequestLogger
// It adds log as the logger and will not log any route in skipRoutes.
func HTTPGenericMiddleware(log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := logmiddleware.HTTPRequestLogger(skipRoutes)(next)
			h = logmiddleware.HTTPAddLogger(log)(h)
			h = trackmiddleware.TrackingID(h)
			h.ServeHTTP(w, r)
		})
	}
}

func TrackingID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackingID := ExtractTrackingID(r)

		ctx := tracking.SetContextID(r.Context(), trackingID)
		rr := r.WithContext(ctx)

		next.ServeHTTP(w, rr)
	})
}

// RequestID is deprecated, use TrackingID instead.
// Deprecated.
func RequestID(next http.Handler) http.Handler {
	return TrackingID(next)
}

func ExtractTrackingID(r *http.Request) string {
	trackingID := r.Header.Get(constants.HeaderTrackingID)
	if trackingID != "" {
		return trackingID
	}
	requestID := r.Header.Get(constants.HeaderRequestID)
	if requestID != "" {
		return requestID
	}
	return uuid.New().String()
}
