package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blacklane/go-libs/x/events"
	"github.com/caarlos0/env"
	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	loadEnvVars()
	log.Printf("config: %#v", cfg)

	errHandler := func(event events.Event, err error) {
		log.Panicf("failed to deliver the event %s: %v", string(event.Payload), err)
	}

	kpc := events.NewKafkaProducerConfig(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaBootstrapServer,
		"message.timeout.ms": 6000,
	})
	kpc.WithEventDeliveryErrHandler(errHandler)

	p, err := events.NewKafkaProducer(kpc)
	if err != nil {
		log.Panicf("could not create kafka producer: %v", err)
	}

	// handle failed deliveries
	_ = p.HandleEvents()
	defer func() {
		if err := p.Shutdown(context.Background()); err != nil {
			log.Panicf("could not shutdown kafka producer: %v", err)
		}
		log.Printf("Shutdown.")
	}()

	e := events.Event{Payload: []byte("Hello, Gophers")}

	err = p.Send(e, cfg.Topic)
	if err != nil {
		log.Panicf("error sending the event %s: %v", e.TopicPartition.Topic, err)
	}
	log.Printf("sent event.Payload: %v", string(e.Payload))
}
