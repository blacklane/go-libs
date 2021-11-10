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
	defer func() { log.Printf("Shutdown: %v", p.Shutdown(context.Background())) }()

	payload := fmt.Sprintf("[%s] Hello, Gophers", time.Now())
	e := events.Event{Payload: []byte(payload)}

	err = p.Send(e, topic)
	if err != nil {
		log.Printf("[ERROR] sending the event %s: %v", e, err)
	}
	log.Printf("publishing: %s", payload)
}
