package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/x/events"
)

//revive:disable:exported Obvious config struct.
type config struct {
	KafkaBootstrapServer string `env:"KAFKA_BOOTSTRAP_SERVERS,required"`
	KafkaGroupID         string `env:"KAFKA_GROUP_ID,required"`
	Topic                string `env:"KAFKA_TOPIC" envDefault:"test-go-libs-events"`
}

var cfg config

func loadEnvVars() {
	c := &config{}
	if err := env.Parse(c); err != nil {
		panic(fmt.Sprintf("could not load environment variables: %v", err))
	}

	cfg = *c
}
func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	loadEnvVars()
	log.Printf("config: %#v", cfg)

	conf := &kafka.ConfigMap{
		"group.id":           cfg.KafkaGroupID,
		"bootstrap.servers":  "localhost:9092",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	c, err := events.NewKafkaConsumer(
		events.NewKafkaConsumerConfig(conf),
		[]string{cfg.Topic},
		events.HandlerFunc(
			func(ctx context.Context, e events.Event) error {
				log.Printf("consumed event: %s", e.Payload)
				return nil
			}))
	if err != nil {
		panic(err)
	}

	log.Printf("starting to consume events from %s, press CTRL+C to exit", cfg.Topic)
	c.Run(time.Second)

	<-signalChan
	if err = c.Shutdown(context.Background()); err != nil {
		log.Printf("Shutdown error: %v", err)
		os.Exit(1)
	}
	log.Printf("Shutdown successfully")
}
