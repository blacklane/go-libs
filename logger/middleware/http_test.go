package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

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
	want := `{"level":"info","application":"TestHTTPRequestLogger","host":"example.com","ip":"localhost","params":"","path":"/do_not_skip","request_id":"","tracking_id":"","user_agent":"","verb":"GET","http_status":200,"duration_ms":0,"body":null,"timestamp":"2009-11-10T23:00:00.000Z","message":"GET /do_not_skip"}` + "\n"

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
		"http://example.com/do_not_skip", nil)
	rLog.Header.Set(internal.HeaderForwardedFor, "localhost")
	rLog.Header.Set(internal.HeaderRequestID, "a_known_id")

	log := logger.New(buf, "TestHTTPRequestLogger")
	ms := chain{HTTPAddLogger(log), HTTPRequestLogger(skipRoutes)}
	h := ms.apply(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(``))
	}))

	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rSkip)
	h.ServeHTTP(w, rLog)

	got := buf.String()
	if want != got {
		t.Errorf("\n\twant: %s,\n\tgot: %s\n\tdiff: %s",
			want, got, cmp.Diff(want, got))
	}
}

func TestHTTPAddBodyFilters(t *testing.T) {
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})
	type args struct {
		filterKeys []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "working with body",
			args: args{
				filterKeys: []string{"hello"},
			},
			want: `{"level":"info","application":"TestHTTPRequestLogger","host":"example.com","ip":"192.0.2.1","params":"","path":"/with_body","request_id":"","tracking_id":"","user_agent":"","verb":"POST","http_status":200,"duration_ms":0,"body":{"hello":"[FILTERED]"},"timestamp":"2009-11-10T23:00:00.000Z","message":"POST /with_body"}` + "\n",
		},
		{
			name: "working body with embedded object",
			args: args{
				filterKeys: []string{"hello", "bye"},
			},
			want: `{"level":"info","application":"TestHTTPRequestLogger","host":"example.com","ip":"192.0.2.1","params":"","path":"/with_body","request_id":"","tracking_id":"","user_agent":"","verb":"POST","http_status":200,"duration_ms":0,"body":{"another":{"bye":"[FILTERED]"},"hello":"[FILTERED]"},"timestamp":"2009-11-10T23:00:00.000Z","message":"POST /with_body"}` + "\n",
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
			log := logger.New(buf, "TestHTTPRequestLogger")
			ms := chain{HTTPAddLogger(log), HTTPAddBodyFilters(tt.args.filterKeys), HTTPRequestLogger([]string{})}
			h := ms.apply(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(``))
			}))

			h.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.want, buf.String())

		case "working body with embedded object":
			buf := bytes.NewBufferString("")
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"http://example.com/with_body", bytes.NewBuffer([]byte(`{"hello":"world","another":{"bye":"bye bye"}}`)))
			log := logger.New(buf, "TestHTTPRequestLogger")
			ms := chain{HTTPAddLogger(log), HTTPAddBodyFilters(tt.args.filterKeys), HTTPRequestLogger([]string{})}
			h := ms.apply(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(``))
			}))

			h.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.want, buf.String())
		}
	}
}
