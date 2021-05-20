package main

import (
	"context"
	"os"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opentracing/opentracing-go"

	"github.com/blacklane/go-libs/tracing"
)

func eventsHandler() events.HandlerFunc {
	return func(ctx context.Context, e events.Event) error {
		logger.FromContext(ctx).Info().Msgf("consumed event: %s", e.Payload)
		return nil
	}
}

func newConsumer(serviceName string, tracer opentracing.Tracer, topic string) events.Consumer {
	// Creates a logger for this "service"
	log := logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName)

	conf := &kafka.ConfigMap{
		"group.id":           "consumer-example",
		"bootstrap.servers":  "localhost:9092",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	// Add the opentracing middleware which parses a span from the headers and
	// injects it on the context. If not span is found, it creates one.
	handler := tracing.EventsAddDefault(eventsHandler(), log, tracer, eventName)

	c, err := events.NewKafkaConsumer(
		events.NewKafkaConsumerConfig(conf),
		[]string{topic},
		handler)
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("starting to consume events from topic: %s", topic)
	c.Run(time.Second)

	return c
}

func newProducer() events.Producer {
	errHandler := func(event events.Event, err error) {
		log.Err(err).Msgf("failed to deliver the event %s",
			string(event.Payload))
	}

	kpc := events.NewKafkaProducerConfig(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"message.timeout.ms": 1000,
	})
	kpc.WithEventDeliveryErrHandler(errHandler)

	p, err := events.NewKafkaProducer(kpc)
	if err != nil {
		log.Panic().Msgf("could not create kafka producer: %v", err)
	}

	// handle failed deliveries
	_ = p.HandleEvents()

	return p
}
