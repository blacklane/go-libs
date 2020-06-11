package events

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func createTopic(topic string, t *testing.T) {
	t.Helper()

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		t.Fatalf("failed to create Admin client: %s", err)
	}
	defer adminClient.Close()

	maxDuration := 60 * time.Second
	results, err := adminClient.CreateTopics(context.Background(),
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1}},
		kafka.SetAdminOperationTimeout(maxDuration))

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
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		t.Fatalf("failed to create producer: %v", err)
	}

	topic := "demo_topic"
	createTopic(topic, t)

	messages := []string{"Anderson", "likes", "reviewing!"}
	handler := ErrorHandlerFunc(func(event Event, err error) {
		message := string(event.Payload)
		t.Errorf("failed to deliver the event %s: %v", message, err)
	})

	p := NewKafkaProducer(kafkaProducer, handler)
	p.HandleMessages()
	defer p.Shutdown(context.Background())

	for _, message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e, topic); err != nil {
			t.Errorf("error sending the event %s: %v", e, err)
		}
	}
}

func TestProducerProducesEventsToIncorrectTopicWithError(t *testing.T) {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		t.Fatalf("failed to create producer: %v", err)
		return
	}

	topic := "not_created_topic"
	messages := map[string]bool{"Anderson": true, "likes": true, "reviewing!": true}
	mu := sync.Mutex{}

	handler := ErrorHandlerFunc(func(event Event, err error) {
		message := string(event.Payload)

		mu.Lock()
		defer mu.Unlock()
		delete(messages, message)
	})

	p := NewKafkaProducer(kafkaProducer, handler)
	p.HandleMessages()
	for message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e, topic); err != nil {
			t.Errorf("error sending the event %s: %v", e, err)
		}
	}

	p.Shutdown(context.Background())

	mu.Lock()
	defer mu.Unlock()
	if len(messages) != 0 {
		t.Errorf("error handler missed %d messages", len(messages))
	}
}
