package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"

	"github.com/blacklane/go-libs/otel"

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

	// OpenTelemetry (OTel) Jaeger tracer
	otel.SetUpOTel(serviceName, cfg.OTelExporterEndpoint, cfg.Log, otel.WithServiceVersion(serviceVersion))

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
	log.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}
