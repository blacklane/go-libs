package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/blacklane/go-libs/tracking"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
)

func ExampleHTTPAddLogger() {
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	r.Header.Set(internal.HeaderForwardedFor, "localhost")

	log := logger.New(prettyWriter, "ExampleHTTPAddLogger")

	loggerMiddleware := HTTPAddLogger(log)

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			log.Info().Msg("Hello, Gophers")
		}))

	h.ServeHTTP(w, r)

	// Output:
	// {
	//   "application": "ExampleHTTPAddLogger",
	//   "level": "info",
	//   "message": "Hello, Gophers",
	//   "timestamp": "2009-11-10T23:00:00Z"
	// }
}

func ExampleHTTPRequestLogger_simple() {
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	r.Header.Set(internal.HeaderForwardedFor, "localhost")

	log := logger.New(prettyWriter, "ExampleHTTPRequestLogger_simple")
	ctx := log.WithContext(r.Context())

	loggerMiddleware := HTTPRequestLogger([]string{})

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "Hello, world")
		}))
	h.ServeHTTP(w, r.WithContext(ctx))

	// Output:
	// {
	//   "application": "ExampleHTTPRequestLogger_simple",
	//   "duration_ms": 0,
	//   "host": "example.com",
	//   "http_status": 200,
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "GET /foo",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "",
	//   "timestamp": "2009-11-10T23:00:00Z",
	//   "tracking_id": "",
	//   "user_agent": "",
	//   "verb": "GET"
	// }

}

func ExampleHTTPRequestLogger_complete() {
	sec := -1
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		sec++
		return time.Date(2009, time.November, 10, 23, 0, sec, 0, time.UTC)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	r.RemoteAddr = "42.42.42.42:42"

	ctx := tracking.SetContextID(r.Context(), "42")

	log := logger.New(prettyWriter, "ExampleHTTPRequestLogger_complete")
	ctx = log.WithContext(ctx)

	rr := r.WithContext(ctx)
	loggerMiddleware := HTTPRequestLogger([]string{})

	h := loggerMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write(nil) }))

	h.ServeHTTP(w, rr)

	// Output:
	// {
	//   "application": "ExampleHTTPRequestLogger_complete",
	//   "duration_ms": 1000,
	//   "host": "example.com",
	//   "http_status": 200,
	//   "ip": "42.42.42.42",
	//   "level": "info",
	//   "message": "GET /foo",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "42",
	//   "timestamp": "2009-11-10T23:00:02Z",
	//   "tracking_id": "42",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
}

func ExampleHTTPRequestLogger_skipRoutes() {
	livePath := "/live"
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	w := httptest.NewRecorder()
	rBar := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	rBar.Header.Set(internal.HeaderForwardedFor, "localhost")

	rLive := httptest.NewRequest(http.MethodGet, "http://example.com"+livePath, nil)
	rLive.Header.Set(internal.HeaderForwardedFor, "localhost")

	log := logger.New(prettyWriter, "ExampleHTTPRequestLogger_skipRoutes")
	ctx := log.WithContext(rBar.Context())

	loggerMiddleware := HTTPRequestLogger([]string{livePath})

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "Hello, world")
		}))

	h.ServeHTTP(w, rBar.WithContext(ctx))
	h.ServeHTTP(w, rLive.WithContext(ctx))

	// Output:
	// {
	//   "application": "ExampleHTTPRequestLogger_skipRoutes",
	//   "duration_ms": 0,
	//   "host": "example.com",
	//   "http_status": 200,
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "GET /foo",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "",
	//   "timestamp": "2009-11-10T23:00:00Z",
	//   "tracking_id": "",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
}

func ExampleHTTPAddAll() {
	trackingID := "tracking_id_ExampleHTTPAddAll"
	livePath := "/live"
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	log := logger.New(prettyWriter, "ExampleHTTPAddAll")

	respWriterBar := httptest.NewRecorder()
	requestBar := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	requestBar.Header.Set(internal.HeaderTrackingID, trackingID) // This header is set to have predictable value in the log output
	requestBar.Header.Set(internal.HeaderForwardedFor, "localhost")

	respWriterLive := httptest.NewRecorder()
	requestLive := httptest.NewRequest(http.MethodGet, "http://example.com"+livePath, nil)
	requestLive.Header.Set(internal.HeaderForwardedFor, "localhost")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "ExampleHTTPAddAll")
	})

	allInOneMiddleware := HTTPAddAll(log, []string{livePath})
	h := allInOneMiddleware(handler)

	h.ServeHTTP(respWriterLive, requestLive)
	h.ServeHTTP(respWriterBar, requestBar)

	// Output:
	// {
	//   "application": "ExampleHTTPAddAll",
	//   "duration_ms": 0,
	//   "host": "example.com",
	//   "http_status": 200,
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "GET /foo",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "tracking_id_ExampleHTTPAddAll",
	//   "timestamp": "2009-11-10T23:00:00Z",
	//   "tracking_id": "tracking_id_ExampleHTTPAddAll",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
}
