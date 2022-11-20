package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/otel"
	"github.com/blacklane/go-libs/x/events"
	"github.com/blacklane/go-libs/x/events/consumer"
	"github.com/blacklane/go-libs/x/events/examples/oauth"
)

var config oauth.Config

func loadEnvVars() {
	c := &oauth.Config{}
	if err := env.Parse(c); err != nil {
		panic(fmt.Sprintf("could not load environment variables: %v", err))
	}

	config = *c
}

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	errLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	loadEnvVars()

	tokenSource := events.NewTokenSource(
		config.ClientID,
		config.ClientSecret,
		config.OauthTokenURL,
		5*time.Second,
		http.Client{Timeout: 3 * time.Second})

	kafkaConfig := &kafka.ConfigMap{
		"group.id":           config.KafkaGroupID,
		"bootstrap.servers":  config.KafkaServer,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	kc := events.NewKafkaConsumerConfig(kafkaConfig)
	kc.WithOAuth(tokenSource)
	kc.WithErrFunc(func(err error) { errLogger.Print(err) })

	log.Printf("creating kafka consumer for topic %s...", config.Topic)
	c, err := events.NewKafkaConsumer(
		kc,
		[]string{config.Topic},
		// using kafka event handler
		events.HandlerFunc(func(ctx context.Context, e events.Event) error {
			log.Printf("consumed event: %s", e.Payload)
			return nil
		}),
		// using consumer for specific event with otel
		consumer.New("my-event", func(ctx context.Context, m consumer.Message) error {
			var payload struct {
				ExpectedProp string `json:"expectedProp"`
			}
			if err := m.DecodePayload(&payload); err != nil {
				return err
			}
			log.Printf("consumed event with expected prop: %s", payload.ExpectedProp)
			return nil
		}, otel.EventConsumer()),
	)
	if err != nil {
		panic(fmt.Sprintf("could not create kafka consumer: %v", err))
	}

	log.Printf("starting consumer listening to %s, press CTRL+C to exit", config.Topic)
	c.Run(time.Second)

	<-signalChan
	if err = c.Shutdown(context.Background()); err != nil {
		log.Printf("Shutdown error: %v", err)
		os.Exit(1)
	}
	log.Printf("Shutdown successfully")
}
