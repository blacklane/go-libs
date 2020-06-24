// +build integration

package events

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestKafkaConsumer_Run(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	mu := sync.Mutex{}

	payloads := map[string]bool{
		fmt.Sprint(time.Now().Unix()) + "-Hello":   true,
		fmt.Sprint(time.Now().Unix()) + "-Gophers": true}
	topic := "TestKafkaConsumer_Run"

	createTopic(t, topic)
	producer := newProducer(t)

	// using a unique groupID to avoid losing events for other consumers
	groupID := "TestKafkaConsumer_Run-" + fmt.Sprint(time.Now().Unix())
	config := &kafka.ConfigMap{
		"group.id":          groupID,
		"bootstrap.servers": kafkaBootstrapServers,
		// for unknown reasons any values smaller than 6000 produces the error:
		// "JoinGroup failed: Broker: Invalid session timeout"
		"session.timeout.ms":       6000,
		"go.events.channel.enable": true,
		"auto.offset.reset":        "earliest",
	}
	c, _ := NewKafkaConsumer(config, []string{topic}, HandlerFunc(func(ctx context.Context, e Event) error {
		mu.Lock()
		defer mu.Unlock()
		delete(payloads, string(e.Payload))

		return nil
	}))

	for msg := range payloads {
		produce(t, producer, msg, topic)
	}

	c.Run(time.Second)

	// We need wait a bit for the messages to get published and consumed
	time.Sleep(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := c.Shutdown(ctx); err != nil {
		t.Errorf("consumer shutdown failed: %v", err)
	}

	if len(payloads) != 0 {
		for p := range payloads {
			t.Errorf(`event "%s" was not consumed`, p)
		}
	}
}

func TestKafkaConsumer_ShutdownTimesOut(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	want := ErrShutdownTimeout

	topic := "TestKafkaConsumer_ShutdownTimesOut"

	createTopic(t, topic)

	// using a unique groupID to avoid losing events for other consumers
	groupID := "TestKafkaConsumer_ShutdownTimesOut-" + fmt.Sprint(time.Now().Unix())
	config := &kafka.ConfigMap{
		"group.id":          groupID,
		"bootstrap.servers": kafkaBootstrapServers,
		// for unknown reasons any values smaller than 6000 produces the error:
		// "JoinGroup failed: Broker: Invalid session timeout"
		"session.timeout.ms":       6000,
		"go.events.channel.enable": true,
		"auto.offset.reset":        "earliest",
	}
	c, _ := NewKafkaConsumer(config, []string{topic}, HandlerFunc(func(context.Context, Event) error { return nil }))

	c.Run(-1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got := c.Shutdown(ctx)

	if !errors.Is(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
