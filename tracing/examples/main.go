package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/logger"
	"github.com/caarlos0/env"

	"github.com/blacklane/go-libs/tracing/internal/jeager"
)

const serviceNameA = "a-service"
const serviceNameB = "b-service"

var eventName = "tracing-example-event"
var topic = "tracing-example"
var tracerHost string

type config struct {
	KafkaServer  string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"http://localhost:9092"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID" envDefault:"tracing-example"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"tracing-example"`
	TracerHost   string `env:"TRACER_HOST" envDefault:"http://localhost:14268/api/traces"`
}

var cfg = config{}
var log = logger.New(logger.ConsoleWriter{Out: os.Stdout}, "tracing-example")

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("failed to load environment variables")
	}

	topic = cfg.Topic
	tracerHost = cfg.TracerHost
}

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	// Opentracing tracer for service A.
	// In production it'll be DataDog's tracer: https://pkg.go.dev/gopkg.in/DataDog/dd-trace-go.v1/ddtrace#example-package-Datadog
	aTracer, aCloser := jeager.NewTracer(tracerHost, serviceNameA, log)
	defer aCloser.Close()

	// Opentracing tracer for service B
	// In production it'll be DataDog's tracer: https://pkg.go.dev/gopkg.in/DataDog/dd-trace-go.v1/ddtrace#example-package-Datadog
	bTracer, bCloser := jeager.NewTracer(tracerHost, serviceNameB, log)
	defer bCloser.Close()

	// Simulates a service consuming messages from Kafka, listening the topic defined by the `topic` variable.
	c := newConsumer(serviceNameB, bTracer, topic)
	defer c.Shutdown(context.TODO())

	p := newProducer()
	defer p.Shutdown(context.TODO())

	// Simulates a service providing a HTTP API on localhost:4242.
	newHTTPServer(serviceNameA, aTracer, p)

	<-signalChan
	log.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}
