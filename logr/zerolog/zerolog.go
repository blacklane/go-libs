package zerolog

import (
	"github.com/blacklane/go-libs/logr"
	"github.com/blacklane/go-libs/logr/field"
	"github.com/rs/zerolog"
)

func New(logger zerolog.Logger) logr.Logger {
	return &zlogr{&logger}
}

type zlogr struct {
	L *zerolog.Logger
}

func (z *zlogr) Info(msg string, fields ...field.Field) {
	eventWithFields(z.L.Info(), fields...).Msg(msg)
}

func (z *zlogr) Error(err error, msg string, fields ...field.Field) {
	eventWithFields(z.L.Err(err), fields...).Msg(msg)
}

func (z *zlogr) WithFields(fields ...field.Field) logr.Logger {
	logger := contextWithFields(z.L.With(), fields...).Logger()
	return &zlogr{&logger}
}
