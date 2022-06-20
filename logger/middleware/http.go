package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
	"github.com/blacklane/go-libs/tracking"
	"github.com/rs/zerolog"
)

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

// HTTPAddBodyFilters adds body filter values into the request context.
func HTTPAddBodyFilters(filterKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), internal.FilterKeys, filterKeys)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// HTTPRequestLogger produces a log line with the request status and fields
// following Blacklane's logging standard.
func HTTPRequestLogger(skipRoutes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := logger.Now()
			urlPath := strings.Split(r.URL.Path, "?")[0] // TODO: obfuscate query string values and show the keys
			ctx := r.Context()
			log := *logger.FromContext(ctx)

			trackingID := tracking.IDFromContext(ctx)
			var body []byte

			//save body to log later
			if r.Body != http.NoBody {
				var err error
				body, err = ioutil.ReadAll(r.Body) //Body swap
				if err != nil {
					log.Log().Msg("error reading body")
				}
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}

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

				var b map[string]interface{}
				if body != nil {
					keys := getKeys(ctx, log)
					b = filterBody(body, keys)
				}

				//post process fields
				f := make(map[string]interface{})
				f[internal.FieldBody] = b

				getLogLevel(log, ww).
					Int(internal.FieldHTTPStatus, ww.statusCode).
					Dur(internal.FieldRequestDuration, logger.Now().Sub(startTime)).
					Fields(f).
					Msgf("%s %s", r.Method, urlPath)
			}()

			next.ServeHTTP(&ww, r)
		})
	}
}

func getKeys(ctx context.Context, log logger.Logger) []string {
	keys, ok := ctx.Value(internal.FilterKeys).([]string)
	if !ok {
		log.Log().Msg("couldn't get filter keys from context")
	}

	if len(keys) == 0 {
		return internal.DefaultKeys
	}

	return keys
}

func filterBody(body []byte, filterKeys []string) map[string]interface{} {
	var b map[string]interface{}
	err := json.Unmarshal(body, &b)
	if err != nil {
		return nil
	}

	for k, v := range b {
		if _, ok := v.(map[string]interface{}); ok {
			marshalData, err := json.Marshal(v)
			if err != nil {
				return nil
			}
			b[k] = filterBody(marshalData, filterKeys)
		}

		//validate tag list
		toSearch := strings.Join(filterKeys, ",")
		if strings.Contains(toSearch, k) {
			b[k] = internal.FilterTag
		}
	}
	return b
}

// Logger adds the logger into the request context.
// Deprecated, use HTTPAddLogger instead.
func Logger(log logger.Logger) func(http.Handler) http.Handler {
	return HTTPAddLogger(log)
}

// RequestLogger use HTTPRequestLogger instead.
// Deprecated
func RequestLogger(skipRoutes []string) func(http.Handler) http.Handler {
	return HTTPRequestLogger(skipRoutes)
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
