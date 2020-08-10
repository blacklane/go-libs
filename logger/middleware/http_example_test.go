package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

	log := logger.New(os.Stdout, "")

	loggerMiddleware := HTTPAddLogger(log)

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			log.Info().Msg("Hello, Gophers")
		}))

	h.ServeHTTP(w, r)

	// Output:
	// {"level":"info","application":"","timestamp":"2009-11-10T23:00:00Z","message":"Hello, Gophers"}
}

func ExampleHTTPRequestLogger_simple() {
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	r.Header.Set(internal.HeaderForwardedFor, "localhost")

	log := logger.New(os.Stdout, "")
	ctx := log.WithContext(r.Context())

	loggerMiddleware := HTTPRequestLogger([]string{})

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "Hello, world")
		}))
	h.ServeHTTP(w, r.WithContext(ctx))

	// Output:
	// {"level":"info","application":"","entry_point":true,"host":"example.com","ip":"localhost","params":"bar=foo","path":"/foo","request_depth":0,"request_id":"","route":"","tree_path":"","user_agent":"","verb":"GET","event":"request_finished","status":200,"request_duration":0,"timestamp":"2009-11-10T23:00:00Z","message":"GET /foo"}
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

	log := logger.New(os.Stdout, "")
	ctx = log.WithContext(ctx)

	rr := r.WithContext(ctx)
	loggerMiddleware := HTTPRequestLogger([]string{})

	h := loggerMiddleware(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) { _, _ = w.Write(nil) }))

	h.ServeHTTP(w, rr)

	// Output:
	// {"level":"info","application":"","entry_point":true,"host":"example.com","ip":"42.42.42.42","params":"bar=foo","path":"/foo","request_depth":0,"request_id":"42","route":"","tree_path":"","user_agent":"","verb":"GET","event":"request_finished","status":200,"request_duration":1000,"timestamp":"2009-11-10T23:00:02Z","message":"GET /foo"}
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

	log := logger.New(os.Stdout, "")
	ctx := log.WithContext(rBar.Context())

	loggerMiddleware := HTTPRequestLogger([]string{livePath})

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "Hello, world")
		}))

	h.ServeHTTP(w, rBar.WithContext(ctx))
	h.ServeHTTP(w, rLive.WithContext(ctx))

	// Output:
	// {"level":"info","application":"","entry_point":true,"host":"example.com","ip":"localhost","params":"bar=foo","path":"/foo","request_depth":0,"request_id":"","route":"","tree_path":"","user_agent":"","verb":"GET","event":"request_finished","status":200,"request_duration":0,"timestamp":"2009-11-10T23:00:00Z","message":"GET /foo"}
}
