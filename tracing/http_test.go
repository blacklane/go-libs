package tracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blacklane/go-libs/tracking"
	"github.com/opentracing/opentracing-go"

	"github.com/blacklane/go-libs/tracing/internal/jeager"
)

func TestHTTPAddOpentracing_noSpan(t *testing.T) {
	wantTrackingID := "TestHTTPAddOpentracing_noSpan"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateRootSpan(t, SpanFromContext(r.Context()), wantTrackingID)
	})

	tracer, closer := jeager.NewTracer("", "Opentracing-integration-test", noopLogger)
	defer func() {
		if err := closer.Close(); err != nil {
			t.Errorf("error closing tracer: %v", err)
		}
	}()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(tracking.SetContextID(r.Context(), wantTrackingID))

	h := HTTPAddOpentracing("/TestHTTPAddOpentracing_noSpan", tracer, handler)
	h.ServeHTTP(httptest.NewRecorder(), r)
}

func TestHTTPAddOpentracing_existingSpan(t *testing.T) {
	wantTrackingID := "TestHTTPAddOpentracing_noSpan"
	path := "/TestHTTPAddOpentracing_noSpan"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateChildSpan(t, SpanFromContext(r.Context()), wantTrackingID)
	})

	tracer, closer := jeager.NewTracer("", "Opentracing-integration-test", noopLogger)
	defer func() {
		if err := closer.Close(); err != nil {
			t.Errorf("error closing tracer: %v", err)
		}
	}()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(tracking.SetContextID(r.Context(), wantTrackingID))

	span := tracer.StartSpan(path)
	err := tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		t.Fatalf("could not inject span into event headers: %v", err)
	}

	h := HTTPAddOpentracing(path, tracer, handler)
	h.ServeHTTP(httptest.NewRecorder(), r)
}
