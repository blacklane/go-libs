package otel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blacklane/go-libs/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type (
	// Option applies a configuration to config.
	Option func(config *Config)

	// Config holds the OTel configuration and is edited by Option.
	Config struct {
		debug            bool
		env              string
		exporterEndpoint string
		serviceName      string
		serviceVersion   string
		errHandler       otel.ErrorHandler
	}
)

// String returns a JSON representation of c. If json.Marshal fails,
// the returned string will be the error.
func (c Config) String() string {
	bs, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("could not marshal otel config to print it: %v", err)
	}

	return fmt.Sprintf(`%s`, bs)
}

// WithServiceVersion adds version as the service version span attribute.
func WithServiceVersion(version string) Option {
	return func(cfg *Config) {
		cfg.serviceVersion = version
	}
}

// WithEnvironment adds env as the environment span attribute.
func WithEnvironment(env string) Option {
	return func(cfg *Config) {
		cfg.env = env
	}
}

// WithDebug enables debug by adding a span processor which prints to stdout.
func WithDebug() Option {
	return func(cfg *Config) {
		cfg.debug = true
	}
}

// WithErrorHandler registers h as OTel's global error handler.
// See go.opentelemetry.io/otel.ErrorHandler for more details on the error handler.
func WithErrorHandler(h func(error)) Option {
	return func(c *Config) {
		c.errHandler = otel.ErrorHandlerFunc(h)
	}
}

// SetUpOTel perform all necessary initialisations for open telemetry and registers
// a trace provider. Any call to OTel API before the setup is done, will likely
// use the default noop implementations.
// serviceName cannot be empty, it identifies the service being instrumented.
// exporterEndpoint is to where the traces will be sent to using GRPC. If it's
// empty, OTel will NOT be enabled, no tracer will be registered.Therefore, the
// OTel APIs will use the default noop implementation.
// log is a logger used to log relevant as well as debug information. If a non-fatal
// error occurs, it's logged as warning and the setup proceeds.
// Check the WithXxx functions for optional configurations.
func SetUpOTel(serviceName, exporterEndpoint string, log logger.Logger, opts ...Option) error {
	cfg := &Config{
		debug:            false,
		env:              "env not set",
		exporterEndpoint: exporterEndpoint,
		serviceName:      serviceName,
		serviceVersion:   "version not set",
	}
	if cfg.serviceName == "" {
		return errors.New("serviceName cannot be empty")
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.exporterEndpoint == "" {
		log.Info().Msg("otel is disabled as OTEL exporter endpoint is empty")
		return nil
	}

	log.Debug().Str("configuration", cfg.String()).Msg("otel configuration")

	if cfg.errHandler != nil {
		otel.SetErrorHandler(cfg.errHandler)
	}

	otlpClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(exporterEndpoint))

	otlpExporter, err := otlptrace.New(context.TODO(), otlpClient)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create OTel exporter, disabling OTel")
		return nil
	}

	// Create a sdk/resource to decorate the app
	// with common attributes from OTel spec
	res, err := resource.New(context.TODO(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.serviceName),
			semconv.ServiceVersionKey.String(cfg.serviceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.env),
		),
	)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create otel sdk/resource")
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithBatcher(otlpExporter),
	)
	if cfg.debug {
		log.Debug().Msg("adding stdout span processor")

		stdoutExporter, err := stdouttrace.New()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize stdouttrace export pipeline")
		}

		tracerProvider.RegisterSpanProcessor(
			trace.NewSimpleSpanProcessor(stdoutExporter))
	}

	// Register the tracer provider and propagator
	// so libraries and frameworks used in the app
	// can reuse it to generate traces and metrics
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)

	return nil
}