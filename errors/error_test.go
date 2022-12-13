package errors

import "testing"

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
