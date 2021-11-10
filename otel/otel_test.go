package otel

import (
	"testing"

	"go.opentelemetry.io/otel"
)

func TestWithDebug(t *testing.T) {
	cfg := Config{}
	WithDebug()(&cfg)

	if cfg.debug != true {
		t.Error("want config.debug = true, got false")
	}
}

func TestWithEnvironment(t *testing.T) {
	want := "testing"
	cfg := Config{}
	WithEnvironment(want)(&cfg)

	if cfg.env != want {
		t.Errorf("want config.env = %s, got %s", want, cfg.env)
	}
}

func TestWithServiceVersion(t *testing.T) {
	want := "42.0.1"
	cfg := Config{}
	WithServiceVersion(want)(&cfg)

	if cfg.serviceVersion != want {
		t.Errorf("want config.serviceVersion = %s, got %s", want, cfg.env)
	}
}

func TestWithErrorHandler(t *testing.T) {
	want := otel.ErrorHandlerFunc(func(err error) {})

	cfg := Config{}
	WithErrorHandler(want)(&cfg)

	if cfg.errHandler == nil {
		t.Errorf("want config.errHandler set, got nil")
	}
}
