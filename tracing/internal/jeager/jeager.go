package jeager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blacklane/go-libs/logger"
)

// JaegerLogger implements the jaeger.Logger interface for logger.Logger
type JaegerLogger logger.Logger

//revive:disable // Internal package and docs on the type definition.
func (jl JaegerLogger) Error(msg string) {
	l := logger.Logger(jl)
	l.Err(errors.New(msg)).Msg("jeager tracer error")
}

func (jl JaegerLogger) Infof(msg string, args ...interface{}) {
	jl.Debugf(msg, args)
}

func (jl JaegerLogger) Debugf(msg string, args ...interface{}) {
	l := logger.Logger(jl)
	l.Debug().Msgf("%s", fmt.Sprintf(strings.TrimSpace(msg), args...))
}

//revive:enable

//revive:disable-line // Obvious function, unnecessary to document.
