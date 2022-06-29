package logr

import (
	"github.com/blacklane/go-libs/logr/field"
)

type Logger interface {
	Debug(msg string, fields ...field.Field)
	Info(msg string, fields ...field.Field)
	Error(err error, msg string, fields ...field.Field)
	WithFields(fields ...field.Field) Logger
	SkipCallerFrame() Logger
}

var defaultLogger Logger = Discard()

func SetLogger(logger Logger) {
	if logger == nil {
		defaultLogger = Discard()
	} else {
		defaultLogger = logger
	}
}

func Debug(msg string, fields ...field.Field) {
	defaultLogger.SkipCallerFrame().Debug(msg, fields...)
}

func Info(msg string, fields ...field.Field) {
	defaultLogger.SkipCallerFrame().Info(msg, fields...)
}

func Error(err error, msg string, fields ...field.Field) {
	defaultLogger.SkipCallerFrame().Error(err, msg, fields...)
}

func WithFields(fields ...field.Field) Logger {
	return defaultLogger.WithFields(fields...)
}
