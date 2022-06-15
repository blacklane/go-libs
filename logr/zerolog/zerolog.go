package zerolog

import (
	"github.com/blacklane/go-libs/logr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(logger zerolog.Logger) logr.Logger {
	log.Debug().Msg("Logger initialized")
	return &zlogr{&logger}
}

type zlogr struct {
	L *zerolog.Logger
}

func (z *zlogr) Info(msg string, fields ...logr.Field) {
	eventWithFields(z.L.Info(), fields...).Msg(msg)
}

func (z *zlogr) Error(msg string, fields ...logr.Field) {
	eventWithFields(z.L.Error(), fields...).Msg(msg)
}

func (z *zlogr) WithFields(fields ...logr.Field) logr.Logger {
	logger := contextWithFields(z.L.With(), fields...).Logger()
	return &zlogr{&logger}
}
