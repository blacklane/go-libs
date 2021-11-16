package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/blacklane/go-libs/logger"

	"github.com/blacklane/go-libs/middleware"
	"github.com/blacklane/go-libs/middleware/internal/constants"
)

func ExampleHTTP() {
	trackingID := "tracking_id_ExampleHTTP"
	ignoredRoute := "/live"

	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	log := logger.New(prettyJSONWriter{}, "ExampleHTTP")

	respWriterBar := httptest.NewRecorder()
	requestBar := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	requestBar.Header.Set(constants.HeaderTrackingID, trackingID)
	requestBar.Header.Set(constants.HeaderForwardedFor, "localhost")

	respWriterIgnored := httptest.NewRecorder()
	requestIgnored := httptest.NewRequest(http.MethodGet, "http://example.com"+ignoredRoute, nil)
	requestIgnored.Header.Set(constants.HeaderTrackingID, trackingID)
	requestIgnored.Header.Set(constants.HeaderForwardedFor, "localhost")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).Info().Msg("always logged")
		_, _ = fmt.Fprint(w, "ExampleHTTP")
	})

	allInOneMiddleware := middleware.HTTP(
		"ExampleHTTPDefaultMiddleware",
		"foo",
		"/foo",
		log,
		ignoredRoute)
	h := allInOneMiddleware(handler)

	h.ServeHTTP(respWriterIgnored, requestIgnored)
	h.ServeHTTP(respWriterBar, requestBar)

	// Output:
	// {
	//   "application": "ExampleHTTP",
	//   "host": "example.com",
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "always logged",
	//   "params": "",
	//   "path": "/live",
	//   "request_id": "tracking_id_ExampleHTTP",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "tracking_id_ExampleHTTP",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
	// {
	//   "application": "ExampleHTTP",
	//   "host": "example.com",
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "always logged",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "tracking_id_ExampleHTTP",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "tracking_id_ExampleHTTP",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
	// {
	//   "application": "ExampleHTTP",
	//   "duration_ms": 0,
	//   "host": "example.com",
	//   "http_status": 200,
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "GET /foo",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "tracking_id_ExampleHTTP",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "tracking_id_ExampleHTTP",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
}

func ExampleHTTP_onlyRequestIDHeader() {
	trackingID := "ExampleHTTP_onlyRequestIDHeader"

	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	log := logger.New(prettyJSONWriter{}, "ExampleHTTP")

	respWriterBar := httptest.NewRecorder()
	requestBar := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	requestBar.Header.Set(constants.HeaderRequestID, trackingID)
	requestBar.Header.Set(constants.HeaderForwardedFor, "localhost")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).Info().Msg("always logged")
		_, _ = fmt.Fprint(w, "ExampleHTTP")
	})

	allInOneMiddleware := middleware.HTTP(
		"ExampleHTTPDefaultMiddleware",
		"foo",
		"/foo",
		log)
	h := allInOneMiddleware(handler)

	h.ServeHTTP(respWriterBar, requestBar)

	// Output:
	// {
	//   "application": "ExampleHTTP",
	//   "host": "example.com",
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "always logged",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "ExampleHTTP_onlyRequestIDHeader",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleHTTP_onlyRequestIDHeader",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
	// {
	//   "application": "ExampleHTTP",
	//   "duration_ms": 0,
	//   "host": "example.com",
	//   "http_status": 200,
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "GET /foo",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "ExampleHTTP_onlyRequestIDHeader",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleHTTP_onlyRequestIDHeader",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
}

func ExampleHTTP_trackingIDAndRequestIDHeaders() {
	trackingID := "ExampleHTTP_trackingIDAndRequestIDHeaders"

	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	log := logger.New(prettyJSONWriter{}, "ExampleHTTP")

	respWriterBar := httptest.NewRecorder()
	requestBar := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	requestBar.Header.Set(constants.HeaderTrackingID, trackingID)
	requestBar.Header.Set(constants.HeaderRequestID, "should ignore")
	requestBar.Header.Set(constants.HeaderForwardedFor, "localhost")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).Info().Msg("always logged")
		_, _ = fmt.Fprint(w, "ExampleHTTP")
	})

	allInOneMiddleware := middleware.HTTP(
		"ExampleHTTPDefaultMiddleware",
		"foo",
		"/foo",
		log)
	h := allInOneMiddleware(handler)

	h.ServeHTTP(respWriterBar, requestBar)

	// Output:
	// {
	//   "application": "ExampleHTTP",
	//   "host": "example.com",
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "always logged",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "ExampleHTTP_trackingIDAndRequestIDHeaders",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleHTTP_trackingIDAndRequestIDHeaders",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
	// {
	//   "application": "ExampleHTTP",
	//   "duration_ms": 0,
	//   "host": "example.com",
	//   "http_status": 200,
	//   "ip": "localhost",
	//   "level": "info",
	//   "message": "GET /foo",
	//   "params": "bar=foo",
	//   "path": "/foo",
	//   "request_id": "ExampleHTTP_trackingIDAndRequestIDHeaders",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "ExampleHTTP_trackingIDAndRequestIDHeaders",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
}
