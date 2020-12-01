package events

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var createdTopics []string
var kafkaBootstrapServers string
var kafkaAdminClient *kafka.AdminClient

func TestMain(m *testing.M) {
	// FYI: go test does not parse any flag when TestMain is defined
	flag.Parse()

	setUp()

	exitCode := m.Run()

	cleanUp()

	os.Exit(exitCode)
}

func setUp() {
	bootstrapServers, ok := os.LookupEnv("KAFKA_BOOTSTRAP_SERVERS")
	if !ok {
		log.Printf("KAFKA_BOOTSTRAP_SERVERS not set, using localhost:9092")
		bootstrapServers = "localhost:9092"
	}
	kafkaBootstrapServers = bootstrapServers

	kad, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBootstrapServers})
	if err != nil {
		panic(fmt.Sprintf("failed to create Admin client: %s", err))
	}

	kafkaAdminClient = kad
}

func cleanUp() {
	if len(createdTopics) > 0 {
		cleanUpTopics()
	}

	kafkaAdminClient.Close()
}

// createTopic creates a topic using the default kafkaAdminClient and adds the
// topic for deletion when all tests finish running.
func createTopic(t *testing.T, topic string) {
	t.Helper()

	results, err := kafkaAdminClient.CreateTopics(context.Background(),
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1}},
		// TODO: double check which ones are really needed
		kafka.SetAdminOperationTimeout(10*time.Second),
		kafka.SetAdminRequestTimeout(10*time.Second))
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}

	createdTopics = append(createdTopics, topic)

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError &&
			result.Error.Code() != kafka.ErrTopicAlreadyExists {
			t.Fatalf("topic creation failed for %s: %s",
				result.Topic, result.Error.String())
		}
	}
}

func cleanUpTopics() {
	results, err :=
		kafkaAdminClient.DeleteTopics(context.Background(), createdTopics)
	if err != nil {
		log.Printf("[ERROR] failed to create topic: %v", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			log.Printf("[ERROR] topic deletion failed for %s: %s",
				result.Topic, result.Error.String())
		}
	}
}

func newConsumer(t *testing.T, topic string, groupID string) *kafka.Consumer {
	t.Helper()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"group.id":          groupID,
		"bootstrap.servers": kafkaBootstrapServers,
		// for unknown reasons any values smaller than 6000 produces the error:
		// "JoinGroup failed: Broker: Invalid session timeout"
		"session.timeout.ms":       6000,
		"go.events.channel.enable": true,
		"auto.offset.reset":        "earliest",
	})
	if err != nil {
		t.Fatalf("could not create kafka consumer: %v", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		t.Fatalf("failed to subscribe to kafka topic %s: %v", topic, err)
	}

	return c
}

func newProducer(t *testing.T) *kafka.Producer {
	t.Helper()

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaBootstrapServers,
		"message.timeout.ms": 1000})
	if err != nil {
		t.Fatalf("failed to create kafkaProducer: %v", err)
	}

	return p
}

func produce(t *testing.T, p *kafka.Producer, key, msg, topic string) {
	t.Helper()

	e := Event{Key: []byte(key),Payload: []byte(msg)}

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny},
		Key:   e.Key,
		Value: e.Payload,
	}, nil)
	if err != nil {
		t.Errorf("error procucing kafka message %v: %v", e, err)
	}
}
