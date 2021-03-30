package logger

import (
	"errors"
	"fmt"
	"io"

	"github.com/rs/zerolog"

	"github.com/blacklane/go-libs/logger/internal"
)

type config struct {
	fields map[string]string
	level  string
}

type Logger = zerolog.Logger
type ConsoleWriter = zerolog.ConsoleWriter

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

// New creates a new logger writing to w.
// Use WithStr(key, value) to add new fields to the logger, example:
//   New(os.Stdout, "myAwesomeApp", WithStr("foo", "bar"))
func New(w io.Writer, appName string, configs ...func(cfg *config)) Logger {
	zerolog.TimestampFieldName = internal.FieldTimestamp
	zerolog.TimeFieldFormat = internal.TimeFieldFormat

	c := zerolog.New(w).
		With().
		Timestamp().
		Str(internal.FieldApplication, appName)

	cfg := newConfig()
	for _, f := range configs {
		f(cfg)
	}

	for k, v := range cfg.fields {
		c = c.Str(k, v)
	}

	l := c.Logger()

	if cfg.level != "" {
		if zl, err := parseLevel(cfg.level); err == nil {
			l = l.Level(zl)
		}
	}

	return l
}

// WithStr adds a new field to the logger with the given key and value.
func WithStr(key, value string) func(cfg *config) {
	return func(cfg *config) {
		cfg.fields[key] = value
	}
}

// WithLevel sets the minimum level to be logged.
// The log levels are: DebugLevel, InfoLevel, WarnLevel and ErrorLevel.
// If an invalid level is given it'll have no effect.
// Use ParseLevel to ensure the level is valid
func WithLevel(level string) func(cfg *config) {
	return func(cfg *config) {
		cfg.level = level
	}
}

// ParseLevel parses the given level, if the level is invalid
// it returns empty string and an error, the given level and nil otherwise
func ParseLevel(level string) (string, error) {
	_, err := parseLevel(level)
	if err != nil {
		return "", err
	}

	return level, nil
}

func newConfig() *config {
	return &config{
		fields: map[string]string{},
	}
}

func parseLevel(level string) (zerolog.Level, error) {
	if level == "" {
		return zerolog.NoLevel,
			errors.New("empty string is not a valid log level")
	}

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return logLevel, fmt.Errorf(
			"%s is not a valid log level: %w", level, err)
	}
	return logLevel, nil
}
