package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/blacklane/go-libs/x/events"
)

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	topic := "events-example-topic"

	conf := &kafka.ConfigMap{
		"group.id":           "consumer-example" + strconv.Itoa(int(time.Now().Unix())),
		"bootstrap.servers":  "localhost:9092",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	c, _ := events.NewKafkaConsumer(conf, []string{topic}, events.HandlerFunc(
		func(ctx context.Context, e events.Event) error {
			log.Printf("consumed event: %s", e.Payload)
			return nil
		}))

	log.Print("starting consumer, press CTRL+C to exit")
	c.Run(time.Second)

	<-signalChan
	log.Printf("Shutdown: %v", c.Shutdown(context.Background()))
}
