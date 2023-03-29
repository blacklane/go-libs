package logger

import (
	"os"
	"strings"
)

type LogOutput string

const (
	LogOutputJSON    LogOutput = "json"
	LogOutputConsole LogOutput = "console"
)

func ParseLogOutput(output string) LogOutput {
	switch strings.ToLower(output) {
	case string(LogOutputConsole):
		return LogOutputConsole
	default:
		return LogOutputJSON
	}
}

type config struct {
	fields map[string]string
	level  string
	env    string
	output LogOutput
}

type Option func(*config)

func defaultConfig() *config {
	cfg := &config{
		fields: map[string]string{},
		output: LogOutputJSON,
		env:    "development",
	}

	if v := os.Getenv("LOG_LEVEL"); v != "" {
		WithLevel(v)(cfg)
	}

	if v := os.Getenv("LOG_OUTPUT"); v != "" {
		WithOutput(ParseLogOutput(v))(cfg)
	}

	if v := os.Getenv("ENV"); v != "" {
		WithEnv(v)(cfg)
	}

	return cfg
}

// WithStr adds a new field to the logger with the given key and value.
func WithStr(key, value string) Option {
	return func(cfg *config) {
		cfg.fields[key] = value
	}
}

// WithLevel sets the minimum level to be logged.
// The log levels are: debug, info, warn and error.
// If an invalid level is given it'll have no effect.
// Use ParseLevel to ensure the level is valid
func WithLevel(level string) Option {
	return func(cfg *config) {
		cfg.level = level
	}
}

// WithOutput sets the logger output format to be used.
func WithOutput(output LogOutput) Option {
	return func(cfg *config) {
		cfg.output = output
	}
}

// WithEnv sets the logger env value to be logged.
func WithEnv(env string) Option {
	return func(cfg *config) {
		cfg.env = env
	}
}
