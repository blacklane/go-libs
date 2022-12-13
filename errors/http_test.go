package errors

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSend(t *testing.T) {
	tests := []struct {
		err  error
		code int
		body string
	}{
		{
			err:  New(503, "my_message", "my_detail"),
			code: 503,
			body: `{"errors":[{"code":503,"message":"my_message","detail":"my_detail"}]}`,
		},
		{
			err:  NotFound("not_found", "not found"),
			code: 404,
			body: `{"errors":[{"code":404,"message":"not_found","detail":"not found"}]}`,
		},
		{
			err:  InternalServer("internal_error", "").WithCause(NotFound("not_found", "not found")),
			code: 500,
			body: `{"errors":[{"code":500,"message":"internal_error"},{"code":404,"message":"not_found","detail":"not found"}]}`,
		},
		{
			err:  errors.New("generic error"),
			code: 500,
			body: `{"errors":[{"code":500,"message":"unknown_error","detail":"generic error"}]}`,
		},
		{
			err:  fmt.Errorf("something happens: %w", New(404, "not_found", "")),
			code: 404,
			body: `{"errors":[{"code":404,"message":"not_found"}]}`,
		},
		{
			err:  fmt.Errorf("something happens: %w", New(404, "not_found", "").WithCause(InternalServer("database_crash", ""))),
			code: 404,
			body: `{"errors":[{"code":404,"message":"not_found"},{"code":500,"message":"database_crash"}]}`,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			r := httptest.NewRecorder()
			Send(r, test.err)

			if r.Code != test.code {
				t.Errorf("invalid code, want: %v, got: %v", test.code, r.Code)
			}

			if body := strings.TrimRight(r.Body.String(), "\n"); body != test.body {
				t.Errorf("invalid body, want: %v, got: %v", test.body, body)
			}
		})
	}
}
