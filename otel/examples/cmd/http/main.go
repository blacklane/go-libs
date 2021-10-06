package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"

	"github.com/blacklane/go-libs/otel"

	"github.com/blacklane/go-libs/otel/examples"
)

const eventName = "event-tracing-example"
const serviceName = "http-example"
const serviceVersion = "devel"

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	cfg := examples.ParseConfig(serviceName)

	// OpenTelemetry (OTel) tracer for service A.
	err := otel.SetUpOTel(
		serviceName,
		cfg.OTelExporterEndpoint,
		cfg.Log,
		otel.WithServiceVersion(serviceVersion),
		otel.WithDebug(),
		otel.WithErrorHandler(func(err error) {
			cfg.Log.Err(err).
				Str("error_type", fmt.Sprintf("%T", err)).
				Msg("an otel irremediable error happened")
		}))
	if err != nil {
		cfg.Log.Panic().Err(err).Msg("SetUpOTel returned and error")
	}

	p := examples.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaServer,
		"message.timeout.ms": 1000,
	}, cfg.Log)
	defer p.Shutdown(context.TODO())

	// Simulates a service providing a HTTP API on localhost:4242.
	examples.StartHTTPServer(serviceName, p, cfg.Topic, eventName)

	<-signalChan
	log.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}
