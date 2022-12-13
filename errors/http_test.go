package errors

import (
	"bytes"
	"errors"
	"fmt"
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
			body: `{"errors":[{"code":503,"message":"my_message","detail":"my_detail"}]}`,
		},
		{
			err:  NotFound("not_found", "not found"),
			body: `{"errors":[{"code":404,"message":"not_found","detail":"not found"}]}`,
		},
		{
			err:  InternalServer("internal_error", "").WithCause(NotFound("not_found", "not found")),
			body: `{"errors":[{"code":500,"message":"internal_error"},{"code":404,"message":"not_found","detail":"not found"}]}`,
		},
		{
			err:  errors.New("generic error"),
			body: `{"errors":[]}`,
		},
		{
			err:  fmt.Errorf("something happens: %w", New(404, "not_found", "")),
			body: `{"errors":[{"code":404,"message":"not_found"}]}`,
		},
		{
			err:  fmt.Errorf("something happens: %w", New(404, "not_found", "").WithCause(InternalServer("database_crash", ""))),
			body: `{"errors":[{"code":404,"message":"not_found"},{"code":500,"message":"database_crash"}]}`,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			buff := &bytes.Buffer{}
			WriteJSON(buff, test.err)

			if body := strings.TrimRight(buff.String(), "\n"); body != test.body {
				t.Errorf("invalid body, want: %v, got: %v", test.body, body)
			}
		})
	}
}
