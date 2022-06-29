package logr

import "github.com/blacklane/go-libs/logr/field"

type discardLogger struct{}

func (d *discardLogger) Debug(msg string, fields ...field.Field)            {}
func (d *discardLogger) Info(msg string, fields ...field.Field)             {}
func (d *discardLogger) Error(err error, msg string, fields ...field.Field) {}
func (d *discardLogger) WithFields(fields ...field.Field) Logger            { return d }
func (d *discardLogger) SkipCallerFrame() Logger                            { return d }

func Discard() Logger {
	return &discardLogger{}
}
