package errors

import (
	"errors"
	"fmt"
	"testing"
)

func FuzzNew(f *testing.F) {
	f.Add(10, "message_1", "detail_1")
	f.Add(700, "error", "")
	f.Add(-90, "my_error", "my details")

	f.Fuzz(func(t *testing.T, code int, message string, detail string) {
		err := New(code, message, detail)

		if err.Code != code {
			t.Errorf("invalid error code, want: %v, got: %v", code, err.Code)
		}

		if err.Message != message {
			t.Errorf("invalid error message, want: %v, got: %v", message, err.Message)
		}

		if err.Detail != detail {
			t.Errorf("invalid error detail, want: %v, got: %v", detail, err.Detail)
		}
	})
}

func TestFrom(t *testing.T) {
	e := New(-1, "message", "")

	tests := []struct {
		err      error
		expected *Error
	}{
		{
			err:      e,
			expected: e,
		},
		{
			err:      fmt.Errorf("wrapped error: %w", e),
			expected: e,
		},
		{
			err:      errors.New("dummy error"),
			expected: New(DefaultErrorCode, DefaultErrorMessage, ""),
		},
		{
			err:      nil,
			expected: nil,
		},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			t.Parallel()

			err := From(test.err)

			if !errors.Is(err, test.expected) {
				t.Errorf("invalid error, want: %v, got: %v", test.expected, err)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	e1 := errors.New("error-1")
	e2 := New(-1, "error-2", "")
	e3 := New(-1, "error-3", "")

	tests := []struct {
		err      *Error
		expected error
	}{
		{
			err:      e2.WithCause(e1),
			expected: e1,
		},
		{
			err:      e3.WithCause(e1).WithCause(e2),
			expected: e2,
		},
		{
			err:      e2,
			expected: nil,
		},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			t.Parallel()

			err := test.err.Unwrap()

			if err != test.expected {
				t.Errorf("invalid error, want: %v, got: %v", test.expected, err)
			}
		})
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name   string
		err1   *Error
		err2   error
		equals bool
	}{
		{
			name:   "same code and message with different detail",
			err1:   New(-1, "message", ""),
			err2:   New(-1, "message", "detail"),
			equals: true,
		},
		{
			name:   "same code but different message",
			err1:   New(-1, "message-1", ""),
			err2:   New(-1, "message-2", ""),
			equals: false,
		},
		{
			name:   "different code with same message",
			err1:   New(1, "message", ""),
			err2:   New(2, "message", ""),
			equals: false,
		},
		{
			name:   "same code and message and detail",
			err1:   New(500, "message", ""),
			err2:   New(500, "message", ""),
			equals: true,
		},
		{
			name:   "different error types",
			err1:   New(500, "message", ""),
			err2:   errors.New("dummy error"),
			equals: false,
		},
		{
			name:   "same code and message and detail with different cause",
			err1:   New(500, "message", ""),
			err2:   New(500, "message", "").WithCause(errors.New("inner error")),
			equals: true,
		},
		{
			name: "same code and message and detail with different metadata",
			err1: New(500, "message", ""),
			err2: New(500, "message", "").WithMetadata(map[string]string{
				"user_id": "123",
			}),
			equals: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			equals := test.err1.Is(test.err2)

			if equals != test.equals {
				t.Errorf("want: %v, got: %v", test.equals, equals)
			}
		})
	}
}

func TestWithDetailf(t *testing.T) {
	err := New(-1, "message", "detail")

	tests := []struct {
		err    *Error
		detail string
	}{
		{
			err:    err.WithDetailf("demo"),
			detail: "demo",
		},
		{
			err:    err.WithDetailf("demo %d", 1),
			detail: "demo 1",
		},
		{
			err:    err.WithDetailf("%d %s %d", 5, "det", 3),
			detail: "5 det 3",
		},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			t.Parallel()

			if test.err.Detail != test.detail {
				t.Errorf("invalid error detail, want: %v, got: %v", test.detail, test.err.Detail)
			}
		})
	}
}

func TestWithMetadataw(t *testing.T) {
	err := New(-1, "message", "detail")

	tests := []struct {
		err      *Error
		metadata map[string]string
	}{
		{
			err: err.WithMetadataw(
				"user_id", "12",
			),
			metadata: map[string]string{
				"user_id": "12",
			},
		},
		{
			err: err.WithMetadataw(
				"a", "1",
				"b",
			),
			metadata: map[string]string{
				"a": "1",
				"b": "",
			},
		},
		{
			err:      err.WithMetadataw(""),
			metadata: map[string]string{"": ""},
		},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			t.Parallel()

			if got, want := len(test.err.Metadata), len(test.metadata); got != want {
				t.Errorf("invalid metadata size, want: %d, got: %d", want, got)
			}

			for key, want := range test.metadata {
				if got := test.err.Metadata[key]; got != want {
					t.Errorf("invalid metadata value for key: %s, want: %s, got: %s", key, want, got)
				}
			}
		})
	}
}
