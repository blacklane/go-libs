package events

import (
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConfigOption func(kafka.ConfigMap)

// NewKafkaConfig creates a kafka config object according to confluentic documentation:
// https://docs.confluent.io/platform/current/installation/configuration
func NewKafkaConfig(options ...KafkaConfigOption) *kafka.ConfigMap {
	conf := kafka.ConfigMap{}
	for _, opt := range options {
		opt(conf)
	}
	return &conf
}

// WithCommaSeparatedBootstrapServers sets the list of host/port pairs to use
// for establishing the initial connection to the Kafka cluster.
func WithCommaSeparatedBootstrapServers(servers string) KafkaConfigOption {
	return func(conf kafka.ConfigMap) {
		conf["bootstrap.servers"] = servers
	}
}

// WithBootstrapServers sets the list of host/port pairs to use for establishing
// the initial connection to the Kafka cluster.
func WithBootstrapServers(servers []string) KafkaConfigOption {
	return WithCommaSeparatedBootstrapServers(strings.Join(servers, ","))
}

// WithSessionTimeout sets the timeout used to detect client failures when using Kafka's group management facility.
// The client sends periodic heartbeats to indicate its liveness to the broker.
func WithSessionTimeout(timeout time.Duration) KafkaConfigOption {
	return func(conf kafka.ConfigMap) {
		conf["session.timeout.ms"] = int(timeout.Milliseconds())
	}
}

func WithTopicMetadataRefreshInterval(interval time.Duration) KafkaConfigOption {
	return func(conf kafka.ConfigMap) {
		conf["topic.metadata.refresh.interval.ms"] = int(interval.Milliseconds())
	}
}

func WithLogConnectionClose(logClose bool) KafkaConfigOption {
	return func(conf kafka.ConfigMap) {
		conf["log.connection.close"] = logClose
	}
}

// WithGroupID sets a unique string that identifies the consumer group this consumer belongs to.
func WithGroupID(groupID string) KafkaConfigOption {
	return func(conf kafka.ConfigMap) {
		conf["group.id"] = groupID
	}
}

type OffsetReset string

const (
	// OffsetResetEarliest automatically reset the offset to the earliest offset
	OffsetResetEarliest OffsetReset = "earliest"
	// OffsetResetLatest automatically reset the offset to the latest offset
	OffsetResetLatest OffsetReset = "latest"
	// OffsetResetNone throw exception to the consumer if no previous offset is found for the consumer's group
	OffsetResetNone OffsetReset = "none"
)

// WithAutoOffsetReset specify what to do when there is no initial offset in Kafka or if the current
// offset does not exist any more on the server (e.g. because that data has been deleted)
func WithAutoOffsetReset(offsetReset OffsetReset) KafkaConfigOption {
	return func(conf kafka.ConfigMap) {
		conf["auto.offset.reset"] = string(offsetReset)
	}
}

func WithKeyValue(key string, value interface{}) KafkaConfigOption {
	return func(conf kafka.ConfigMap) {
		conf[key] = value
	}
}
