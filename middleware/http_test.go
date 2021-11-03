package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/blacklane/go-libs/tracking"

	"github.com/blacklane/go-libs/middleware/internal/constants"
)

func TestHTTPAddTracking(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{name: "tracking id in the context",
			ctx:  tracking.SetContextID(context.Background(), "TestHTTPAddTracking"),
			want: "TestHTTPAddTracking"},
		{name: "context without tracking id",
			ctx:  context.Background(),
			want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/ignore", nil)
			if err != nil {
				t.Fatalf("could not create a request: %s", err)
			}
			HTTPAddTrackingIDToRequest(tt.ctx, r)
			got := r.Header.Get(constants.HeaderTrackingID)
			if got != tt.want {
				t.Errorf("got %s, want: %s", got, tt.want)
			}
		})
	}
}
