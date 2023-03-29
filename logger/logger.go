package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/blacklane/go-libs/logger/internal"
)

func init() {
	zerolog.TimestampFieldName = internal.FieldTimestamp
	zerolog.TimeFieldFormat = internal.TimeFieldFormat
}

// Type alias for zerolog types.
//
//revive:disable:exported The whole block is already documented.
type (
	Logger        = zerolog.Logger
	ConsoleWriter = zerolog.ConsoleWriter
)

// NewStdout creates a new logger writing to standard output.
func NewStdout(appName string, options ...Option) Logger {
	return New(os.Stdout, appName, options...)
}

// New creates a new logger writing to w.
// Use WithStr(key, value) to add new fields to the logger, example:
//
//	New(os.Stdout, "myAwesomeApp", WithStr("foo", "bar"))
func New(w io.Writer, appName string, options ...Option) Logger {
	cfg := defaultConfig()
	for _, opt := range options {
		opt(cfg)
	}

	w = logWriter(cfg.output, w)

	c := zerolog.New(w).
		With().
		Timestamp().
		Str(internal.FieldApplication, appName)

	if cfg.env != "" {
		c = c.Str(internal.FieldEnv, cfg.env)
	}

	for k, v := range cfg.fields {
		c = c.Str(k, v)
	}

	l := c.Logger()

	if cfg.level != "" {
		if zl, err := parseLevel(cfg.level); err == nil {
			l = l.Level(zl)
		} else {
			l.Warn().Err(err).Msg("invalid log level")
		}
	}

	return l
}

func logWriter(output LogOutput, w io.Writer) io.Writer {
	switch output {
	case LogOutputConsole:
		return zerolog.ConsoleWriter{Out: w}
	default:
		return w
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
