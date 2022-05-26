package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/otel"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/otel/examples"
)

const eventName = "event-tracing-example"
const serviceName = "events-example"
const serviceVersion = "devel"

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	cfg := examples.ParseConfig(serviceName)

	// Set up OpenTelemetry (OTel) with a Jaeger tracer.
	otel.SetUpOTel(serviceName, cfg.Log,
		otel.WithGrpcTraceExporter(cfg.OTelExporterEndpoint),
		otel.WithDebug(),
		otel.WithServiceVersion(serviceVersion),
		otel.WithErrorHandler(func(err error) {
			cfg.Log.Err(err).
				Str("error_type", fmt.Sprintf("%T", err)).
				Msg("an otel irremediable error happened")
		}))

	kafkaCfg := &kafka.ConfigMap{
		"group.id":           cfg.KafkaGroupID,
		"bootstrap.servers":  cfg.KafkaServer,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	c := examples.NewStartedConsumer(
		serviceName,
		kafkaCfg,
		cfg.Topic,
		eventName)
	defer c.Shutdown(context.TODO())

	<-signalChan
	cfg.Log.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}
