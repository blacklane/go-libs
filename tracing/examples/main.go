package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/logger"
	"github.com/caarlos0/env"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/tracing/internal/jeager"
)

const serviceNameA = "a-service"
const serviceNameB = "b-service"

var eventName = "tracing-example-event"
var topic = "tracing-example"
var tracerHost string

type config struct {
	KafkaServer  string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"localhost:9092"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID" envDefault:"tracing-example"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"tracing-example"`
	TracerHost   string `env:"TRACER_HOST" envDefault:"localhost:55680"`
	// TracerHost   string `env:"TRACER_HOST" envDefault:"http://localhost:14268/api/traces"`
}

var cfg = config{}
var log = logger.New(logger.ConsoleWriter{Out: os.Stdout}, "tracing-example")

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("failed to load environment variables")
	}

	topic = cfg.Topic
	tracerHost = cfg.TracerHost

	log.Debug().Msgf("config: %#v", cfg)
}

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	// OpenTelemetry (OTel) tracer for service A.
	jeager.InitOpenTelemetry("tracing-example-app", "devel", tracerHost)

	// Simulates a service consuming messages from Kafka, listening the topic defined by the `topic` variable.
	c := newConsumer(
		serviceNameB,
		topic,
		&kafka.ConfigMap{
			"group.id":           cfg.KafkaGroupID,
			"bootstrap.servers":  "http://" + cfg.KafkaServer,
			"session.timeout.ms": 6000,
			"auto.offset.reset":  "earliest",
		})
	defer c.Shutdown(context.TODO())

	p := newProducer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaServer,
		"message.timeout.ms": 6000,
	})
	defer p.Shutdown(context.TODO())

	// Simulates a service providing a HTTP API on localhost:4242.
	newHTTPServerWithOTel(serviceNameA, p)

	<-signalChan
	log.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}
