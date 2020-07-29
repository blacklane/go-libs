package middleware

import (
	"context"
	"os"
	"time"

	"github.com/blacklane/go-libs/x/events"

	"github.com/blacklane/go-libs/logger"
)

func ExampleEventsAddLogger() {
	// Set current time function so we can control the logged timestamp
	logger.SetNowFunc(func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	})

	log := logger.New(os.Stdout, "")

	h := events.HandlerFunc(func(ctx context.Context, _ events.Event) error {
		l := logger.FromContext(ctx)
		l.Info().Msg("Hello, Gophers")
		return nil
	})

	m := EventsAddLogger(log)
	hh := m(h)

	_ = hh.Handle(context.Background(), events.Event{})

	// Output:
	// {"level":"info","application":"","timestamp":"2009-11-10T23:00:00Z","message":"Hello, Gophers"}
}
