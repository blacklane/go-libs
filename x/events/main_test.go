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
var timeoutMultiplier = time.Second // a multiplier for timeouts allowing to adjust then to faster or slower environments

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

	if customTimeout, ok := os.LookupEnv("TIMEOUT_MULTIPLIER"); ok {
		timeoutMultiplier, err = time.ParseDuration(customTimeout)
		if err != nil {
			log.Printf("TIMEOUT_MULTIPLIER not set correctly, using fallback - 1 second")
			timeoutMultiplier = time.Second
		}
	}
}

func cleanUp() {
	if len(createdTopics) > 0 {
		cleanUpTopics()
	}

	kafkaAdminClient.Close()
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
