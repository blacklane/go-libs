package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/blacklane/go-libs/logger"
	logmiddleware "github.com/blacklane/go-libs/logger/middleware"
	"github.com/blacklane/go-libs/x/events"
	"github.com/caarlos0/env"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"

	"github.com/blacklane/go-libs/tracking/internal/constants"

	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/tracking/internal/jeager"
	"github.com/blacklane/go-libs/tracking/middleware"
)

const serviceNameA = "a-service"
const serviceNameB = "b-service"

var topic = "tracing-example"
var eventName = "tracing-example-event"

type config struct {
	KafkaServer  string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"http://localhost:9092"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID" envDefault:"tracing-example"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"tracing-example"`
}

var cfg = config{}
var log = logger.New(logger.ConsoleWriter{Out: os.Stdout}, "tracing-example")

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("failed to load environment variables")
	}

	topic = cfg.Topic
}

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	// Opentracing tracer for service A.
	// In production it'll be DataDog's tracer: https://pkg.go.dev/gopkg.in/DataDog/dd-trace-go.v1/ddtrace#example-package-Datadog
	aTracer, aCloser := jeager.NewTracer(serviceNameA, log)
	defer aCloser.Close()

	// Opentracing tracer for service B
	// In production it'll be DataDog's tracer: https://pkg.go.dev/gopkg.in/DataDog/dd-trace-go.v1/ddtrace#example-package-Datadog
	bTracer, bCloser := jeager.NewTracer(serviceNameB, log)
	defer bCloser.Close()

	// Simulates a service consuming messages from kafka.
	// the 'topic' variable.
	c := consumer(serviceNameB, bTracer, topic)
	defer c.Shutdown(context.TODO())

	p := producer()
	defer p.Shutdown(context.TODO())

	// Simulates a service providing a HTTP API on localhost:4242.
	httpServer(serviceNameA, aTracer, p)

	<-signalChan
	log.Info().Msg("Goodbye cruel world!")
	os.Exit(0)
}

func httpServer(serviceName string, tracer opentracing.Tracer, producer events.Producer) {
	// Creates a logger for this "service"
	log := logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName)

	// Add the opentracing middleware which parses a span from the headers and
	// injects it on the context. If not span is found, it creates one.
	_, h := middleware.HTTPAddOpentracing(
		"/tracing/example/path", tracer, httpHandler(producer))

	// Add other middleware to fulfil our tracking standards
	handler := logmiddleware.HTTPAddDefault(log)(h)
	httpServer := http.Server{
		Addr:    ":4242",
		Handler: handler,
	}

	go func() { httpServer.ListenAndServe() }()
	log.Info().Msgf("Starting HTTP server on %s", httpServer.Addr)
}

func eventsHandler() events.HandlerFunc {
	return func(ctx context.Context, e events.Event) error {
		logger.FromContext(ctx).Info().Msgf("consumed event: %s", e.Payload)
		return nil
	}
}

func httpHandler(producer events.Producer) http.Handler {
	var count int
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { count++ }()
		ctx := r.Context()

		// Gets the opentracing span
		sp := tracking.SpanFromContext(ctx)

		// add fields to a log entry on this span
		sp.LogFields(
			opentracinglog.String("some_uuid", uuid.New().String()),
			opentracinglog.String("some_key", "some_value"),
			opentracinglog.Int("some_number", 1500),
			opentracinglog.Message("some message"))

		headers, _ := json.Marshal(w.Header())

		// simulates a failure by flipping a coin.
		if rand.Int()%2 == 0 {
			err := errors.New(http.StatusText(http.StatusTeapot))

			// Flags an error happened on this span
			ext.Error.Set(sp, true)

			// Logs the error on the span
			sp.LogFields(opentracinglog.Error(err))

			w.WriteHeader(http.StatusTeapot)
			_, _ = fmt.Fprintf(w, "I'm a tea pot\ntracking_id: %s\nheaders: %s",
				tracking.IDFromContext(ctx),
				headers)
			return
		}

		event := events.Event{
			Headers: nil,
			Key:     []byte(fmt.Sprintf("%d", count)),
			Payload: []byte(fmt.Sprintf(`{"event":"%s","count":%d}`, eventName, count)),
		}

		// Injects the span into the events headers.
		// It's how the span is propagated through different services.
		err := tracking.EventsOpentracingInject(ctx, sp, &event)
		if err != nil {
			logger.FromContext(ctx).Err(err).
				Str(constants.FieldEvent, eventName).
				Msg("could not inject span into event")
		}

		err = producer.Send(event, topic)
		if err != nil {
			err := fmt.Errorf("could not send event: %w", err)

			// Flags an error happened on this span
			ext.Error.Set(sp, true)

			// Logs the error on the span
			sp.LogFields(opentracinglog.Error(err))

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"error":"internal server error","tracking_id":"%s","headers":%q}`,
				tracking.IDFromContext(ctx),
				headers)
			return
		}

		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(
			fmt.Sprintf(`{"tracking_id":"%s","headers":%q}`,
				tracking.IDFromContext(ctx),
				headers)))
	})
}

func consumer(serviceName string, tracer opentracing.Tracer, topic string) events.Consumer {
	// Creates a logger for this "service"
	log := logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName)

	conf := &kafka.ConfigMap{
		"group.id":           "consumer-example",
		"bootstrap.servers":  "localhost:9092",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	// Add other middleware to fulfil our tracking standards
	handler := logmiddleware.EventsAddDefault(eventsHandler(), log, eventName)

	// Add the opentracing middleware which parses a span from the headers and
	// injects it on the context. If not span is found, it creates one.
	handler = middleware.EventsAddOpentracing(
		eventName, tracer, handler)

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

func producer() events.Producer {
	errHandler := func(event events.Event, err error) {
		log.Panic().Msgf("failed to deliver the event %s: %v",
			string(event.Payload), err)
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
