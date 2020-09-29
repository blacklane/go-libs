package middleware

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/tracking/internal/constants"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ExtractRequestID(r)

		ctx := tracking.SetContextID(r.Context(), requestID)
		rr := r.WithContext(ctx)

		next.ServeHTTP(w, rr)
	})
}

func ExtractRequestID(r *http.Request) string {
	requestID := r.Header.Get(constants.HeaderTrackingID)
	if requestID == "" {
		requestID = r.Header.Get(constants.HeaderRequestID)
	}
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
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
