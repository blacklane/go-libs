package logger

import (
	"time"

	"github.com/rs/zerolog"
)

var now = time.Now

// SetNowFunc will change the time function used by logger.Now(). It defaults to time.Now()
// It will also change zerolog.TimestampFunc accordingly
func SetNowFunc(f func() time.Time) {
	now = f
	zerolog.TimestampFunc = now
}

// Now returns the current time. It defaults to time.Now. To change
// it check SetNowFunc.
func Now() time.Time {
	return now()
}
