package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/x/events"
	"github.com/blacklane/go-libs/x/events/examples/oauth"
	"github.com/blacklane/go-libs/x/events/producer"
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
	errLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	loadEnvVars()
	topic := config.Topic

	errHandler := func(event events.Event, err error) {
		log.Panicf("[ERROR] [ErrorHandler] failed to deliver the event %s: %s",
			string(event.Payload), err)
	}

	tokenSource := events.NewTokenSource(
		config.ClientID,
		config.ClientSecret,
		config.OauthTokenURL,
		5*time.Second,
		http.Client{Timeout: 3 * time.Second})

	kpc := events.NewKafkaProducerConfig(&kafka.ConfigMap{
		"bootstrap.servers":  config.KafkaServer,
		"message.timeout.ms": 10000,
	})
	kpc.WithEventDeliveryErrHandler(errHandler)
	kpc.WithOAuth(tokenSource)
	kpc.WithErrFunc(func(err error) { errLogger.Print(err) })

	log.Printf("creating kafka consumer for topic %s...", topic)
	p, err := events.NewKafkaProducer(kpc)
	if err != nil {
		log.Panicf("could not create kafka producer: %v", err)
	}

	_ = p.HandleEvents()
	defer func() {
		if err := p.Shutdown(context.Background()); err != nil {
			log.Panicf("could not shutdown kafka producer: %v", err)
		}
		log.Printf("Shutdown.")
	}()

	payload := fmt.Sprintf("[%s] Hello, Gophers", time.Now())
	e := events.Event{Payload: []byte(payload)}

	err = p.Send(e, topic)
	if err != nil {
		log.Printf("[ERROR] sending the event %s: %v", e.TopicPartition.Topic, err)
	}
	log.Printf("publishing: %s", payload)

	ctx := context.Background()

	tp := producer.NewTopicProducer(p, topic)
	if err := tp.Send(ctx, "my-event", payload, producer.WithKey("event-key")); err != nil {
		log.Printf("error sending event: %s", err)
	}

	if err := tp.Send(ctx, "my-event", payload); err != nil {
		log.Printf("error sending event: %s", err)
	}
}
