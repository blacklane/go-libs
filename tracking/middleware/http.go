package middleware

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/tracking/internal/constants"
)

func TrackingID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackingID := ExtractTrackingID(r)

		ctx := tracking.SetContextID(r.Context(), trackingID)
		rr := r.WithContext(ctx)

		next.ServeHTTP(w, rr)
	})
}

func ExtractTrackingID(r *http.Request) string {
	requestID := r.Header.Get(constants.HeaderTrackingID)
	if requestID != "" {
		return requestID
	}
	requestID = r.Header.Get(constants.HeaderRequestID)
	if requestID != "" {
		return requestID
	}
	return uuid.New().String()
}

// ExtractRequestDepth returns the request depth extracted from the header added of 1 or
// zero if either no header is found or if any error happens when conventing the header value to int
func ExtractRequestDepth(r *http.Request) int {
	depth, err := strconv.Atoi(r.Header.Get(constants.HeaderRequestDepth))
	if err != nil {
		depth = 0
	} else {
		depth++
	}
	return depth
}

func ExtractTreePath(r *http.Request) string {
	return r.Header.Get(constants.HeaderTreePath)
}
