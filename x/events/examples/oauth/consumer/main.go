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

	"github.com/blacklane/go-libs/x/events"
	"github.com/blacklane/go-libs/x/events/examples/oauth"
)

var config oauth.Cfg

func loadEnvVars() {
	c := &oauth.Cfg{}
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
	topic := config.Topic

	tokenSource := oauth.NewTokenSource(
		config.ClientID,
		config.ClientSecret,
		config.TokenURL,
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

	log.Printf("creating kafka consumer for topic %s...", topic)
	c, err := events.NewKafkaConsumer(
		kc,
		[]string{topic},
		events.HandlerFunc(
			func(ctx context.Context, e events.Event) error {
				log.Printf("consumed event: %s", e.Payload)
				return nil
			}))
	if err != nil {
		panic(fmt.Sprintf("could not create kafka consumer: %v", err))
	}

	log.Printf("starting consumer listening to %s, press CTRL+C to exit", topic)
	c.Run(time.Second)

	<-signalChan
	log.Printf("Shutdown: %v", c.Shutdown(context.Background()))
}
