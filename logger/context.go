package logger

import (
	"context"

	"github.com/blacklane/go-libs/logger/hooks"
	"github.com/rs/zerolog"
)

const (
	ddTraceIDKey = "dd.trace_id"
	ddSpanIDKey  = "dd.span_id"
)

// FromContext delegates to zerolog.Ctx. It'll return the logger in the context
// or a disabled logger if none is found. For details check the docs for zerolog.Ctx.
func FromContext(ctx context.Context) *Logger {
	logger := zerolog.Ctx(ctx)
	logger.Hook(hooks.TraceContext(ctx))
	return logger
}
