package middleware

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/blacklane/go-libs/tracking"
	trackingMiddleware "github.com/blacklane/go-libs/tracking/middleware"
	"github.com/rs/zerolog"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
)

// HTTPAddDefault adds the necessary middleware for:
//   - have tracking id in the context (read from the headers or a new one),
//   - have a logger.Logger with tracking id and all required fields in the context,
//   - log, at the end of handler, if it succeeded or failed and how log it took.
// - github.com/blacklane/go-libs/tracking/middleware.HTTPTrackingID
// - middleware.HTTPAddLogger
// - middleware.HTTPRequestLogger
func HTTPAddDefault(log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
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
				internal.FieldHost:       r.Host,
				internal.FieldIP:         ipAddress(r),
				internal.FieldParams:     r.URL.RawQuery,
				internal.FieldPath:       r.URL.Path,
				internal.FieldRequestID:  trackingID,
				internal.FieldTrackingID: trackingID,
				internal.FieldUserAgent:  r.UserAgent(),
				internal.FieldVerb:       r.Method,
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

				getLogLevel(log, ww).
					Int(internal.FieldHTTPStatus, ww.statusCode).
					Dur(internal.FieldRequestDuration, logger.Now().Sub(startTime)).
					Msgf("%s %s", r.Method, urlPath)
			}()

			next.ServeHTTP(&ww, r)
		})
	}
}

// getLogLevel decides which log level to use, Info, Warn or Error.
func getLogLevel(log logger.Logger, ww responseWriter) *zerolog.Event {
	if ww.StatusCode() >= 400 && ww.StatusCode() < 500 {
		return log.Warn()
	}
	if ww.StatusCode() >= 500 {
		return log.Err(
			fmt.Errorf("request finished with HTTP status %d", ww.StatusCode()))
	}
	return log.Info()
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

// HTTPMagic delegates to HTTPAddDefault
func HTTPMagic(log logger.Logger, skipRoutes ...string) func(http.Handler) http.Handler {
	return HTTPAddDefault(log, skipRoutes...)
}
