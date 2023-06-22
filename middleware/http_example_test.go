package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/middleware"
	"github.com/blacklane/go-libs/middleware/internal/constants"
	"github.com/stretchr/testify/assert"
)

func ExampleHTTP() {
	trackingID := "tracking_id_ExampleHTTP"

	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	log := logger.New(prettyJSONWriter{}, "ExampleHTTP")

	respWriterBar := httptest.NewRecorder()
	requestBar := httptest.NewRequest(http.MethodGet, "http://example.com/foo?bar=foo", nil)
	requestBar.Header.Set(constants.HeaderTrackingID, trackingID)
	requestBar.Header.Set(constants.HeaderForwardedFor, "localhost")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.From(r.Context())
		log.Info().Msg("always logged")
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
	//   "request_id": "tracking_id_ExampleHTTP",
	//   "timestamp": "2009-11-10T23:00:00.000Z",
	//   "trace_id": "00000000000000000000000000000000",
	//   "tracking_id": "tracking_id_ExampleHTTP",
	//   "user_agent": "",
	//   "verb": "GET"
	// }
	// {
	//   "application": "ExampleHTTP",
	//   "body": null,
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
		log := logger.From(r.Context())
		log.Info().Msg("always logged")

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
	//   "body": null,
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
		log := logger.From(r.Context())
		log.Info().Msg("always logged")

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
	//   "body": null,
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

func TestHTTPWithBodyFilter(t *testing.T) {
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	trackingID := "tracking_id_ExampleHTTP"
	type args struct {
		serviceName string
		handlerName string
		path        string
		filterKeys  []string
		log         logger.Logger
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "working with body",
			args: args{
				serviceName: "test",
				handlerName: "test",
				path:        "/test",
				filterKeys:  []string{},
				log:         logger.Logger{},
			},
			want: `{"level":"info","application":"TestHTTPRequestLogger","trace_id":"00000000000000000000000000000000","host":"example.com","ip":"192.0.2.1","params":"","path":"/with_body","request_id":"tracking_id_ExampleHTTP","tracking_id":"tracking_id_ExampleHTTP","user_agent":"","verb":"POST","http_status":200,"duration_ms":0,"body":{"hello":"world"},"timestamp":"2009-11-10T23:00:00.000Z","message":"POST /with_body"}` + "\n",
		},
	}
	for _, tt := range tests {
		switch tt.name {
		case "working with body":
			buf := bytes.NewBufferString("")
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"http://example.com/with_body", bytes.NewBuffer([]byte(`{"hello":"world"}`)))
			r.Header.Set(constants.HeaderTrackingID, trackingID)
			log := logger.New(buf, "TestHTTPRequestLogger")

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(``))
				if err != nil {
					t.Error(err)
				}
			})
			midd := middleware.HTTPWithBodyFilter(tt.args.serviceName, tt.args.handlerName, tt.args.path, tt.args.filterKeys, log)

			h := midd(handler)

			h.ServeHTTP(w, r)
			//assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.want, buf.String())
		}
	}
}
