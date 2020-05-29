package logger

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	errNil := func(err error) bool { return err == nil }
	tcs := []struct {
		name  string
		level string
		errOk func(error) bool
	}{
		{name: "level debug", level: "debug", errOk: errNil},
		{name: "level info", level: "info", errOk: errNil},
		{name: "level warn", level: "warn", errOk: errNil},
		{name: "level error", level: "error", errOk: errNil},
		{
			name:  "level bad",
			level: "",
			errOk: func(err error) bool { return err != nil },
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseLevel(tc.level)
			if !tc.errOk(err) {
				t.Errorf("unexpected err: %v", err)
			}
			if got != tc.level {
				t.Errorf("got: %q, want: %q", got, tc.level)
			}
		})
	}
}

func TestWithStr(t *testing.T) {
	wantKey, wantValue := "wantKey", "wantValue"
	want := newConfig()
	want.fields = map[string]string{wantKey: wantValue}

	got := newConfig()
	WithStr(wantKey, wantValue)(got)

	if len(got.fields) != 1 {
		t.Errorf("want one field, got %d: %v", len(got.fields), got.fields)
	}
	for k, v := range got.fields {
		if k != wantKey {
			t.Errorf("got: %s, want: %s", k, wantKey)
		}
		if v != wantValue {
			t.Errorf("got: %s, want: %s", v, wantValue)
		}
	}
}
