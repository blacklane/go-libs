package zap

import (
	"github.com/blacklane/go-libs/logr"
	"github.com/blacklane/go-libs/logr/field"
	"go.uber.org/zap"
)

func New(logger *zap.Logger) logr.Logger {
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	return &zapLogr{logger}
}

type zapLogr struct {
	L *zap.Logger
}

func (z *zapLogr) Debug(msg string, fields ...field.Field) {
	z.L.Debug(msg, mapFields(fields)...)
}

func (z *zapLogr) Info(msg string, fields ...field.Field) {
	z.L.Info(msg, mapFields(fields)...)
}

func (z *zapLogr) Error(err error, msg string, fields ...field.Field) {
	z.L.Error(msg, mapFields(fields)...)
}

func (z *zapLogr) WithFields(fields ...field.Field) logr.Logger {
	logger := z.L.With(mapFields(fields)...)
	return &zapLogr{logger}
}

func (z *zapLogr) SkipCallerFrame() logr.Logger {
	logger := z.L.WithOptions(zap.AddCallerSkip(1))
	return &zapLogr{logger}
}
