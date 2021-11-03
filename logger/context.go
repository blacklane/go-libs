package logger

import (
	"context"

	"github.com/rs/zerolog"
)

// FromContext delegates to zerolog.Ctx. It'll return the logger in the context
// or a disabled logger if none is found. For details check the docs for zerolog.Ctx.
func FromContext(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}
