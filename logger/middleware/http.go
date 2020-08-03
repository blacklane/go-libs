package middleware

import (
	"bytes"
	"net"
	"net/http"
	"strings"

	"github.com/blacklane/go-libs/tracking"
	"github.com/rs/zerolog"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
)

// HTTPAddLogger adds the logger into the request context.
func HTTPAddLogger(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			l := logger.FromContext(ctx)
			logFields := map[string]interface{}{
				internal.FieldEntryPoint: isEntryPoint(r),
				// TODO double check if they are a must and make them optional if not
				internal.FieldRequestDepth: tracking.RequestDepthFromCtx(ctx),
				internal.FieldRequestID:    tracking.IDFromContext(ctx),
				internal.FieldTreePath:     tracking.TreePathFromCtx(ctx),
				internal.FieldRoute:        tracking.RequestRouteFromCtx(ctx),
				internal.FieldParams:       r.URL.RawQuery,
				internal.FieldIP:           ipAddress(r),
				internal.FieldUserAgent:    r.UserAgent(),
				internal.FieldHost:         r.Host,
				internal.FieldVerb:         r.Method,
				internal.FieldPath:         r.URL.Path,
			}
			l.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Fields(logFields)
			})

			ww := responseWriter{w: w, body: &bytes.Buffer{}}

			defer func() {
				for _, skipRoute := range skipRoutes {
					if strings.HasPrefix(urlPath, skipRoute) {
						// do not log this route
						return
					}
				}
				l.Info().
					Str(internal.FieldEvent, internal.EventRequestFinished).
					Int(internal.FieldStatus, ww.statusCode).
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
	return len(r.Header.Get(internal.HeaderRequestID)) == 0
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
