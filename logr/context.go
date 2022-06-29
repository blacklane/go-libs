package logr

import (
	"context"
	"errors"
)

type contextKey struct{}

var ErrLoggerNotFound = errors.New("logger not found in context")

func FromContext(ctx context.Context) (Logger, error) {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v, nil
	}
	return Discard(), ErrLoggerNotFound
}

func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func FromContextOrDiscard(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}
	return Discard()
}
