// +build integration

package events

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestProducerProducesEventsToCorrectTopic(t *testing.T) {
	topic := "demo_topic"
	createTopic(t, topic)

	messages := []string{"Anderson", "likes", "reviewing!"}
	errHandler := func(event Event, err error) {
		message := string(event.Payload)
		t.Errorf("failed to deliver the event %s: %v", message, err)
	}

	kpc := NewKafkaProducerConfig(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaBootstrapServers,
		"message.timeout.ms": 2000,
	})
	kpc.WithEventDeliveryErrHandler(errHandler)

	p, err := NewKafkaProducer(kpc)
	if err != nil {
		t.Fatalf("%v", err)
	}

	_ = p.HandleEvents()
	defer p.Shutdown(context.Background())

	for _, message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e, topic); err != nil {
			t.Errorf("error sending the event %s: %v", e, err)
		}
	}
}

func TestProducerProducesEventsToIncorrectTopicWithError(t *testing.T) {
	topic := "not_created_topic"
	messages := map[string]bool{"Anderson": true, "likes": true, "reviewing!": true}
	mu := sync.Mutex{}

	errHandler := func(e Event, err error) {
		message := string(e.Payload)
		var errDelivery kafka.Error
		if ok := errors.As(err, &errDelivery); !ok {
			t.Errorf("want %T, got %T", kafka.Error{}, err)
		}
		if errDelivery.Code() != kafka.ErrUnknownTopicOrPart {
			t.Errorf("want %s, got: %s",
				kafka.ErrUnknownTopicOrPart.String(), errDelivery.Code())
		}

		mu.Lock()
		defer mu.Unlock()
		delete(messages, message)
	}

	kpc := NewKafkaProducerConfig(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaBootstrapServers,
		"message.timeout.ms": 2000,
	})
	kpc.WithEventDeliveryErrHandler(errHandler)

	p, err := NewKafkaProducer(kpc)
	if err != nil {
		t.Fatalf("%v", err)
	}

	_ = p.HandleEvents()
	for message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e, topic); err != nil {
			t.Errorf("error sending the event %s: %v", e, err)
		}
	}

	err = p.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("failed Shutdown: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(messages) != 0 {
		t.Errorf("error handler missed %d messages", len(messages))
	}
}

func TestNewKafkaProducerWithFlushTimeout(t *testing.T) {
	want := 1

	kpc := NewKafkaProducerConfig(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaBootstrapServers,
		"message.timeout.ms": 2000,
	})
	kpc.WithFlushTimeout(want)
	p, err := NewKafkaProducer(kpc)
	if err != nil {
		t.Fatalf("could not create a new kafka producer: %v", err)
	}

	kp, ok := p.(*kafkaProducer)
	if !ok {
		t.Errorf("%T is not a %T", p, &kafkaProducer{})
	}

	if kp.flushTimeoutMs != want {
		t.Errorf("got: %d, want: %d", kp.flushTimeoutMs, want)
	}
}
