package middleware

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/blacklane/go-libs/tracking"
	trackingMiddleware "github.com/blacklane/go-libs/tracking/middleware"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
)

// HTTPAddAll
func HTTPAddAll(log logger.Logger, skipRoutes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := HTTPRequestLogger(skipRoutes)(next)
			h = HTTPAddLogger(log)(h)
			h = trackingMiddleware.TrackingID(h)
			h.ServeHTTP(w, r)
		})
	}
}

// HTTPAddLogger adds the logger into the request context.
func HTTPAddLogger(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log = log.With().Logger() // Creates a sub logger so all requests won't share the same logger instance
			ctx := log.WithContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logger adds the logger into the request context.
// Deprecated, use HTTPAddLogger instead.
func Logger(log logger.Logger) func(http.Handler) http.Handler {
	return HTTPAddLogger(log)
}

// HTTPRequestLogger produces a log line with the request status and fields
// following the standard defined on http://handbook.int.blacklane.io/monitoring/kiev.html#requestresponse-logging
func HTTPRequestLogger(skipRoutes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := logger.Now()
			urlPath := strings.Split(r.URL.Path, "?")[0] // TODO: obfuscate query string values and show the keys
			ctx := r.Context()

			log := *logger.FromContext(ctx)
			trackingID := tracking.IDFromContext(ctx)

			logFields := map[string]interface{}{
				internal.FieldTrackingID: trackingID,
				internal.FieldRequestID:  trackingID,
				internal.FieldParams:     r.URL.RawQuery,
				internal.FieldIP:         ipAddress(r),
				internal.FieldUserAgent:  r.UserAgent(),
				internal.FieldHost:       r.Host,
				internal.FieldVerb:       r.Method,
				internal.FieldPath:       r.URL.Path,
			}

			log = log.With().Fields(logFields).Logger()
			ctx = log.WithContext(ctx)
			r = r.WithContext(ctx)

			ww := responseWriter{w: w, body: &bytes.Buffer{}}

			defer func() {
				for _, skipRoute := range skipRoutes {
					if strings.HasPrefix(urlPath, skipRoute) {
						// do not log this route
						return
					}
				}

				zerologEvent := log.Info()
				if ww.StatusCode() >= 400 && ww.StatusCode() < 500 {
					zerologEvent = log.Warn().Err(fmt.Errorf("request finished with http status "))
				}
				zerologEvent.
					Int(internal.FieldHTTPStatus, ww.statusCode).
					Dur(internal.FieldRequestDuration, logger.Now().Sub(startTime)).
					Msgf("%s %s", r.Method, urlPath)
			}()

			next.ServeHTTP(&ww, r)
		})
	}
}

// RequestLogger use HTTPRequestLogger instead.
// Deprecated
func RequestLogger(skipRoutes []string) func(http.Handler) http.Handler {
	return HTTPRequestLogger(skipRoutes)
}

func ipAddress(r *http.Request) string {
	forwardedIP := r.Header.Get(internal.HeaderForwardedFor)
	if len(forwardedIP) == 0 {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		return ip
	}
	return forwardedIP
}

func isEntryPoint(r *http.Request) bool {
	return r.Header.Get(internal.HeaderRequestID) == "" && r.Header.Get(internal.HeaderTrackingID) == ""
}

type responseWriter struct {
	w          http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w *responseWriter) Header() http.Header {
	// TODO: add test
	return w.w.Header()
}

func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

func (w *responseWriter) Body() string {
	// TODO: add test
	w.body.Reset()
	return w.body.String()
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func (w *responseWriter) Write(buf []byte) (int, error) {
	if w.StatusCode() == 0 {
		w.WriteHeader(http.StatusOK)
	}

	if w.statusCode < http.StatusOK || w.statusCode > 299 {
		w.body.Write(buf)
	}

	return w.w.Write(buf)
}
