package otel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blacklane/go-libs/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	ErrEmptyServiceName      = errors.New("serviceName cannot be empty")
	ErrTraceExporterNotFound = errors.New("OTEL Trace Exporter not found")
)

const (
	defaultEnv     = "none"
	defaultVersion = ""
)

type (
	// Option applies a configuration to config.
	Option func(config *Config)

	// Config holds the OTel configuration and is edited by Option.
	Config struct {
		Enable          bool   `json:"enable"`
		Debug           bool   `json:"debug"`
		Env             string `json:"env"`
		otlpTraceClient otlptrace.Client
		ServiceName     string `json:"serviceName"`
		ServiceVersion  string `json:"serviceVersion"`
		errHandler      otel.ErrorHandler
	}
)

func defaultConfig() *Config {
	cfg := &Config{
		Enable:         boolEnv("OTEL_ENABLED", true),
		Debug:          boolEnv("OTEL_DEBUG", false),
		Env:            defaultEnv,
		ServiceName:    filepath.Base(os.Args[0]),
		ServiceVersion: defaultVersion,
	}

	if v := os.Getenv("OTEL_TRACE_EXPORTER_OTLP_HTTP_ENDPOINT"); v != "" {
		WithHttpTraceExporter(v)(cfg)
	}

	if v := os.Getenv("OTEL_TRACE_EXPORTER_OTLP_GRPC_ENDPOINT"); v != "" {
		WithGrpcTraceExporter(v)(cfg)
	}

	if v := os.Getenv("ENV"); v != "" {
		WithEnvironment(v)(cfg)
	}

	if v := os.Getenv("APPLICATION"); v != "" {
		cfg.ServiceName = v
	}

	return cfg
}

func (c *Config) validate() error {
	if c.ServiceName == "" {
		return ErrEmptyServiceName
	}
	if c.otlpTraceClient == nil {
		return ErrTraceExporterNotFound
	}
	return nil
}

// String returns a JSON representation of c. If json.Marshal fails,
// the returned string will be the error.
func (c Config) String() string {
	bs, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("could not marshal otel config to print it: %v", err)
	}

	return string(bs)
}

// WithServiceVersion adds version as the service version span attribute.
func WithServiceVersion(version string) Option {
	return func(cfg *Config) {
		cfg.ServiceVersion = version
	}
}

// WithEnvironment adds env as the environment span attribute.
func WithEnvironment(env string) Option {
	return func(cfg *Config) {
		cfg.Env = env
	}
}

// WithDebug enables debug by adding a span processor which prints to stdout.
func WithDebug() Option {
	return func(cfg *Config) {
		cfg.Debug = true
	}
}

// WithErrorHandler registers h as OTel's global error handler.
// See go.opentelemetry.io/otel.ErrorHandler for more details on the error handler.
func WithErrorHandler(h func(error)) Option {
	return func(c *Config) {
		c.errHandler = otel.ErrorHandlerFunc(h)
	}
}

// WithGrpcTraceExporter registers an otlp trace exporter.
func WithGrpcTraceExporter(endpoint string) Option {
	return func(c *Config) {
		c.otlpTraceClient = otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(endpoint),
		)
	}
}

// WithHttpTraceExporter registers an otlp http trace exporter.
func WithHttpTraceExporter(endpoint string) Option {
	return func(c *Config) {
		c.otlpTraceClient = otlptracehttp.NewClient(
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(endpoint),
		)
	}
}

// SetUpOTel perform all necessary initialisations for open telemetry and registers
// a trace provider. Any call to OTel API before the setup is done, will likely
// use the default noop implementations.
//
// some values will be infered from env vars:
//
// - OTEL_ENABLED: to activate/deactivate otel, default: true
//
// - OTEL_DEBUG: to activate debug mode, default: false
//
// - OTEL_TRACE_EXPORTER_OTLP_HTTP_ENDPOINT: otlp HTTP trace exporter is activated
//
// - OTEL_TRACE_EXPORTER_OTLP_GRPC_ENDPOINT: otlp GRPC trace exporter is activated
//
// - ENV: is the application environment
func SetUpOTel(serviceName string, log logger.Logger, opts ...Option) error {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	if !cfg.Enable {
		log.Info().Msg("otel is disabled")
		return nil
	}

	if cfg.errHandler == nil {
		cfg.errHandler = otel.ErrorHandlerFunc(func(err error) {
			log.Err(err).Msg("otel internal error")
		})
	}
	cfg.ServiceName = serviceName

	if err := cfg.validate(); err != nil {
		log.Err(err).Msg("invalid otel configuration")
		return err
	}

	log.Debug().Str("configuration", cfg.String()).Msg("otel configuration")

	otlpExporter, err := createTraceExporter(cfg)
	if err != nil {
		log.Err(err).Msg("failed to create trace exporter")
		return err
	}

	// Create a sdk/resource to decorate the app
	// with common attributes from OTel spec
	res, err := createResource(cfg)
	if err != nil {
		log.Err(err).Msg("failed to create otel sdk/resource")
		return err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithBatcher(otlpExporter),
	)
	if cfg.Debug {
		log.Debug().Msg("adding stdout span processor")

		stdoutExporter, err := stdouttrace.New()
		if err != nil {
			log.Err(err).Msg("failed to initialize stdouttrace export pipeline")
		}

		tracerProvider.RegisterSpanProcessor(
			trace.NewSimpleSpanProcessor(stdoutExporter))
	}

	// Register the tracer provider and propagator
	// so libraries and frameworks used in the app
	// can reuse it to generate traces and metrics
	otel.SetTracerProvider(tracerProvider)
	otel.SetErrorHandler(cfg.errHandler)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)

	log.Info().Msg("otel started")
	return nil
}

func createTraceExporter(cfg *Config) (*otlptrace.Exporter, error) {
	if cfg.otlpTraceClient == nil {
		return nil, ErrTraceExporterNotFound
	}
	return otlptrace.New(context.TODO(), cfg.otlpTraceClient)
}

func createResource(cfg *Config) (*resource.Resource, error) {
	return resource.New(context.TODO(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Env),
		),
	)
}
