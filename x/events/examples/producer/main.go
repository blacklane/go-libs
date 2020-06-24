package main

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/x/events"
)

func main() {
	errHandler := events.ErrorHandlerFunc(func(event events.Event, err error) {
		log.Panicf("failed to deliver the event %s: %v", string(event.Payload), err)
	})

	p, err := events.NewKafkaProducer(
		&kafka.ConfigMap{
			"bootstrap.servers":  "localhost:9092",
			"message.timeout.ms": "1000"},
		errHandler)
	if err != nil {
		log.Panicf("%v", err)
	}

	// handle failed deliveries
	_ = p.HandleMessages()
	defer p.Shutdown(context.Background())

	e := events.Event{Payload: []byte("Hello, Gophers")}

	err = p.Send(e, "events-example-topic")
	if err != nil {
		log.Panicf("error sending the event %s: %v", e, err)
	}
}
