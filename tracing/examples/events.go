package examples

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/tracing"
)

func eventsHandler() events.HandlerFunc {
	return func(ctx context.Context, e events.Event) error {
		time.Sleep(time.Duration(5*rand.Intn(5)) * time.Millisecond)

		logger.FromContext(ctx).Info().
			Str("event_headers", fmt.Sprintf("%v", e.Headers)).
			Str("event_payload", string(e.Payload)).
			Msg("consumed event")
		return nil
	}
}

// NewStartedConsumer creates, initialises, starts and then returns a new Kafka consumer.
func NewStartedConsumer(serviceName string, conf *kafka.ConfigMap, topic string, eventName string) events.Consumer {
	// Creates a logger for this "service"
	log := logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName).
		With().
		Str("environment", "otel").Logger()

	// Add the opentracing middleware which parses a span from the headers and
	// injects it on the context. If not span is found, it creates one.
	handler := tracing.EventsAddDefault(eventsHandler(), log, eventName)

	c, err := events.NewKafkaConsumer(
		events.NewKafkaConsumerConfig(conf),
		[]string{topic},
		handler)
	if err != nil {
		panic(err)
	}

	c.Run(time.Second)
	log.Info().Msgf("started to consume events from topic: %s", topic)
	log.Info().Msgf("consumer kafka configs: %v", conf)

	return c
}

// NewProducer creates, initialises, starts delivery failure handling and then returns a new Kafka producer.
func NewProducer(conf *kafka.ConfigMap, log logger.Logger) events.Producer {
	errHandler := func(event events.Event, err error) {
		log.Err(err).Msgf("failed to deliver the event %s",
			string(event.Payload))
	}

	kpc := events.NewKafkaProducerConfig(conf)
	kpc.WithEventDeliveryErrHandler(errHandler)

	p, err := events.NewKafkaProducer(kpc)
	if err != nil {
		log.Panic().Msgf("could not create kafka producer: %v", err)
	}

	// handle failed deliveries
	_ = p.HandleEvents()
	log.Info().Msgf("producer kafka configs: %v", conf)

	return p
}
