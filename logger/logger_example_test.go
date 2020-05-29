package logger

import (
	"os"
	"time"
)

func ExampleNew() {
	// It's only needed to have a consistent output for the timestamp
	SetNowFunc(func() time.Time { return time.Time{} })

	l := New(
		os.Stdout,
		"example",
		WithStr("key1", "value1"),
		WithStr("key2", "value2"))

	l.Info().Msg("Hello, Gophers!")

	// Output:
	// {"level":"info","application":"example","key1":"value1","key2":"value2","timestamp":"0001-01-01T00:00:00Z","message":"Hello, Gophers!"}
}
