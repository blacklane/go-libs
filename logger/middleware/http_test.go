package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blacklane/go-libs/tracking/middleware"
	"github.com/google/go-cmp/cmp"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/logger/internal"
)

type chain []func(http.Handler) http.Handler

func (m chain) apply(handler http.Handler) http.Handler {
	h := m[len(m)-1](handler)
	for i := len(m) - 2; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

func TestHTTPRequestLogger(t *testing.T) {
	want := `{"level":"info","application":"TestHTTPRequestLogger","entry_point":false,"host":"example.com","ip":"localhost","params":"","path":"/do_not_skip","request_depth":0,"request_id":"a_known_id","route":"","tree_path":"","user_agent":"","verb":"GET","event":"request_finished","status":0,"request_duration":0,"timestamp":"2009-11-10T23:00:00Z","message":"GET /do_not_skip"}` + "\n"

	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	skipRoutes := []string{"/skip1", "/live"}

	buf := bytes.NewBufferString("")
	w := httptest.NewRecorder()

	rSkip := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://example.com%s", skipRoutes[1]), nil)
	rLog := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://example.com/do_not_skip"), nil)
	rLog.Header.Set(internal.HeaderForwardedFor, "localhost")
	rLog.Header.Set(internal.HeaderRequestID, "a_known_id")

	log := logger.New(buf, "TestHTTPRequestLogger")
	ms := chain{middleware.RequestID, HTTPAddLogger(log), HTTPRequestLogger(skipRoutes)}
	h := ms.apply(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rLog)

	got := buf.String()
	if want != got {
		t.Errorf("\n\twant: %s\tgot: %s\tdiff: %s", want, got, cmp.Diff(want, got))
	}
}
