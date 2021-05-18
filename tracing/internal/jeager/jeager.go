package jeager

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/blacklane/go-libs/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// NewTracer returns a Jeager implementation of opentracing.Tracer
func NewTracer(serviceName string, logger logger.Logger) (opentracing.Tracer, io.Closer) {
	jaegerCfg := jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://localhost:14268/api/traces",
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

func (l JaegerLogger) Error(msg string) {
	log := logger.Logger(l)
	log.Err(errors.New(msg)).Msg("jeager tracer error")
}

func (l JaegerLogger) Infof(msg string, args ...interface{}) {
	l.Debugf(msg, args)
}

func (l JaegerLogger) Debugf(msg string, args ...interface{}) {
	log := logger.Logger(l)
	log.Debug().Msgf("%s", fmt.Sprintf(strings.TrimSpace(msg), args...))
}
