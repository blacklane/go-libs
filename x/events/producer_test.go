package events

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func createTopic(topic string) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		return fmt.Errorf("failed to create Admin client: %s", err)
	}
	defer adminClient.Close()
	maxDuration, err := time.ParseDuration("60s")
	if err != nil {
		return fmt.Errorf("failed to parse time.ParseDuration(60s)")
	}
	results, err := adminClient.CreateTopics(context.Background(),
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1}},
		kafka.SetAdminOperationTimeout(maxDuration))
	if err != nil {
		return fmt.Errorf("failed to create topic: %s", err)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError &&
			result.Error.Code() != kafka.ErrTopicAlreadyExists {
			return fmt.Errorf("topic creation failed for %s: %s",
				result.Topic, result.Error.String())
		}
	}

	return nil
}

func TestProducerProducesEventsToCorrectTopic(t *testing.T) {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		t.Errorf("failed to create producer")
		return
	}
	topic := "demo_topic"
	err = createTopic(topic)
	if err != nil {
		t.Errorf("failed to create topic %s", err)
		return
	}
	messages := []string{"Anderson", "likes", "reviewing!"}
	handler := ErrorHandlerFunc(func(event *Event, err error) {
		if event == nil {
			t.Errorf("failed to parse the event")
			return
		}
		message := string(event.Payload)
		t.Errorf("failed to deliver the event %s", message)
	})
	p := NewKafkaProducer(kafkaProducer, handler)
	p.Run()
	defer p.Shutdown(context.Background())
	for _, message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e, topic); err != nil {
			t.Errorf("error sending the event %s", e)
			return
		}
	}
}

func TestProducerProducesEventsToIncorrectTopicWithError(t *testing.T) {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka"})
	if err != nil {
		t.Errorf("failed to create producer")
		return
	}
	topic := "not_created_topic"
	messages := map[string]bool{"Anderson": true, "likes": true, "reviewing!": true}
	mu := sync.Mutex{}
	handler := ErrorHandlerFunc(func(event *Event, err error) {
		if event == nil {
			t.Errorf("failed to parse the event")
			return
		}
		message := string(event.Payload)
		mu.Lock()
		delete(messages, message)
		mu.Unlock()
	})
	p := NewKafkaProducer(kafkaProducer, handler)
	p.Run()
	for message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e, topic); err != nil {
			t.Errorf("error sending the event %s", e)
		}
	}
	p.Shutdown(context.Background())
	mu.Lock()
	if len(messages) != 0 {
		t.Errorf("error handler missed %d messages", len(messages))
	}
	mu.Unlock()
}
