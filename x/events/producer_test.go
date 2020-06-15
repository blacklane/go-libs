package events

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func createTopic(t *testing.T, topic string) {
	t.Helper()

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		t.Fatalf("failed to create Admin client: %s", err)
	}
	defer adminClient.Close()

	results, err := adminClient.CreateTopics(context.Background(),
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1}},
		kafka.SetAdminOperationTimeout(60*time.Second))
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError &&
			result.Error.Code() != kafka.ErrTopicAlreadyExists {
			t.Fatalf("topic creation failed for %s: %s",
				result.Topic, result.Error.String())
		}
	}
}

func TestProducerProducesEventsToCorrectTopic(t *testing.T) {
	topic := "demo_topic"
	createTopic(t, topic)

	messages := []string{"Anderson", "likes", "reviewing!"}
	handler := ErrorHandlerFunc(func(event Event, err error) {
		message := string(event.Payload)
		t.Errorf("failed to deliver the event %s: %v", message, err)
	})

	p, err := NewKafkaProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"}, handler)
	if err != nil {
		t.Fatalf("%v", err)
	}

	_ = p.HandleMessages()
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

	handler := ErrorHandlerFunc(func(event Event, err error) {
		message := string(event.Payload)

		mu.Lock()
		defer mu.Unlock()
		delete(messages, message)
	})

	p, err := NewKafkaProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"}, handler)
	if err != nil {
		t.Fatalf("%v", err)
	}

	_ = p.HandleMessages()
	for message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e, topic); err != nil {
			t.Errorf("error sending the event %s: %v", e, err)
		}
	}

	_ = p.Shutdown(context.Background())

	mu.Lock()
	defer mu.Unlock()
	if len(messages) != 0 {
		t.Errorf("error handler missed %d messages", len(messages))
	}
}
