package logr

import (
	"github.com/blacklane/go-libs/logr/field"
)

type Logger interface {
	Info(msg string, fields ...field.Field)
	Error(err error, msg string, fields ...field.Field)
	WithFields(fields ...field.Field) Logger
}

type dummyLogger struct{}

func (d *dummyLogger) Info(msg string, fields ...field.Field)             {}
func (d *dummyLogger) Error(err error, msg string, fields ...field.Field) {}
func (d *dummyLogger) WithFields(fields ...field.Field) Logger            { return d }

var defaultLogger Logger = &dummyLogger{}

func SetLogger(logger Logger) {
	if logger == nil {
		defaultLogger = &dummyLogger{}
	} else {
		defaultLogger = logger
	}
}

func Info(msg string, fields ...field.Field) {
	defaultLogger.Info(msg, fields...)
}

func Error(err error, msg string, fields ...field.Field) {
	defaultLogger.Error(err, msg, fields...)
}

func WithFields(fields ...field.Field) Logger {
	return defaultLogger.WithFields(fields...)
}
