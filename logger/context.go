package logger

import (
	"context"

	"github.com/blacklane/go-libs/logger/hooks"
	"github.com/rs/zerolog"
)

// FromContext delegates to zerolog.Ctx. It'll return the logger in the context
// or a disabled logger if none is found. For details check the docs for zerolog.Ctx.
//
// Warning: this function returns a pointer only for backward compatibility, avoid using log.UpdateContext().
//
// Deprecated: This function has been replaced by logger.From(ctx).
func FromContext(ctx context.Context) *Logger {
	logger := zerolog.Ctx(ctx).
		Hook(hooks.TraceContext(ctx))

	return &logger
}

// FromContext delegates to zerolog.Ctx.
// It'll return the logger in the context with hooks for opentelemetry and tracking.
func From(ctx context.Context) Logger {
	return zerolog.Ctx(ctx).
		Hook(hooks.TraceContext(ctx))
}
