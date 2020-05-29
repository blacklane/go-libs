package logger

import (
	"context"

	"github.com/rs/zerolog"
)

func FromContext(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}
