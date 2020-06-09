package events

import (
	"context"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func createTopic(topic string, t *testing.T) {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	ctx, _ := context.WithCancel(context.Background())
	if err != nil {
		t.Errorf("failed to create Admin client: %s\n", err)
		return
	}
	results, err := adminClient.CreateTopics(ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1}},
		kafka.SetAdminOperationTimeout(1000))
	if err != nil {
		t.Errorf("failed to create topic: %s", err)
		return
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError &&
			result.Error.Code() != kafka.ErrTopicAlreadyExists {
			t.Errorf("topic creation failed for %s: %s",
				result.Topic, result.Error.String())
			return
		}
	}

	adminClient.Close()
}

func TestProducerProducesEvents(t *testing.T) {
	configMap := kafka.ConfigMap{"bootstrap.servers": "localhost"}
	kafka.NewAdminClient(&configMap)
	kafkaProducer, err := kafka.NewProducer(&configMap)
	if err != nil {
		t.Errorf("failed to create producer")
	}
	topic := "demo_topic"
	createTopic(topic, t)
	messages := []string{"Anderson", "likes", "reviewing!"}
	handler := ErrorHandlerFunc(func(event *Event, err error) {
		if event == nil {
			t.Errorf("failed to parse the event")
			return
		}
		message := string(event.Payload)
		t.Errorf("failed to deliver the event %s", message)
	})
	p := NewKafkaProducer(kafkaProducer, topic, handler)
	p.Run()
	defer p.Shutdown(context.Background())
	for _, message := range messages {
		e := Event{Payload: []byte(message)}
		if err := p.Send(e); err != nil {
			t.Errorf("error sending the event %s", e)
		}
	}
}
