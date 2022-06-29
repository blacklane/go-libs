package zerolog

import (
	"github.com/blacklane/go-libs/logr"
	"github.com/blacklane/go-libs/logr/field"
	"github.com/rs/zerolog"
)

func New(logger zerolog.Logger) logr.Logger {
	return &zlogr{&logger, 1}
}

type zlogr struct {
	L          *zerolog.Logger
	callerSkip int
}

func (z *zlogr) Debug(msg string, fields ...field.Field) {
	eventWithFields(z.L.Debug().CallerSkipFrame(z.callerSkip), fields...).Msg(msg)
}

func (z *zlogr) Info(msg string, fields ...field.Field) {
	eventWithFields(z.L.Info().CallerSkipFrame(z.callerSkip), fields...).Msg(msg)
}

func (z *zlogr) Error(err error, msg string, fields ...field.Field) {
	eventWithFields(z.L.Err(err).CallerSkipFrame(z.callerSkip), fields...).Msg(msg)
}

func (z *zlogr) WithFields(fields ...field.Field) logr.Logger {
	logger := contextWithFields(z.L.With(), fields...).Logger()
	return &zlogr{&logger, z.callerSkip}
}

func (z *zlogr) SkipCallerFrame() logr.Logger {
	return &zlogr{z.L, z.callerSkip + 1}
}
