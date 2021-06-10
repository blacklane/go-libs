package jeager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/blacklane/go-libs/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

// NewOpentracingTracer returns a Jeager implementation of opentracing.Tracer
// deprecated
func NewOpentracingTracer(host string, serviceName string, logger logger.Logger) (opentracing.Tracer, io.Closer) {
	jaegerCfg := jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: host,
		},
	}

	// Example metrics factory. Use github.com/uber/jaeger-lib/metrics to bind
	// to real metric framework.
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := jaegerCfg.NewTracer(
		jaegerconfig.Logger(
			JaegerLogger(logger.With().
				Str("component", "JaegerTracer").
				Logger())),
		jaegerconfig.Metrics(jMetricsFactory))
	if err != nil {
		panic(fmt.Errorf("could not initialize jaeger tracer: %w", err))
	}

	return tracer, closer
}

// JaegerLogger implements the jaeger.Logger interface for logger.Logger
type JaegerLogger logger.Logger

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

func InitOpenTelemetry(serviceName, serviceVersion, exporterURL string) {
	// Create an gRPC-based OTLP exporter that
	// will receive the created telemetry data
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(exporterURL),
	)
	exporter, err := otlp.NewExporter(context.TODO(), driver)
	if err != nil {
		log.Fatalf("%s: %v", "failed to create exporter", err)
	}

	// Create a sdk/resource to decorate the app
	// with common attributes from OTel spec
	res, err := resource.New(context.TODO(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			semconv.DeploymentEnvironmentKey.String("example"),
		),
	)
	if err != nil {
		log.Fatalf("%s: %v", "failed to create sdk/resource", err)
	}

	// Create a tracer provider that processes
	// spans using a batch-span-processor. This
	// tracer provider will create a sample for
	// every trace created, which is great for
	// demos but horrible for production –– as
	// volume of data generated will be intense
	bsp := tracesdk.NewBatchSpanProcessor(exporter)
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(res),
		tracesdk.WithSpanProcessor(bsp),
	)

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

	return
}
