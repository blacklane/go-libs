package events

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ErrProducerNotHandlingEvents = errors.New("producer should be handling events")
var ErrProducerIsAlreadyRunning = errors.New("producer is already running")

type Producer interface {
	// Send an event to the given topic
	Send(event Event, topic string) error
	// HandleEvents starts to listen to the producer events channel
	HandleEvents() error
	// Shutdown gracefully shuts down the producer, it respect the context
	// timeout.
	Shutdown(ctx context.Context) error
}

// KafkaProducerConfig holds all possible configurations for the kafka producer.
// Use NewKafkaProducerConfig to initialise it.
// To see the possible configurations, check the its WithXXX methods and
// *kafkaConfig.WithXXX` methods as well
type KafkaProducerConfig struct {
	*kafkaConfig

	deliveryErrHandler func(Event, error)

	flushTimeoutMs int
}

type kafkaProducer struct {
	*kafkaCommon

	config   *kafka.ConfigMap
	producer *kafka.Producer

	deliveryErrHandler func(Event, error)

	flushTimeoutMs int

	runMu    sync.RWMutex
	run      bool
	shutdown bool
}

// NewKafkaProducerConfig returns a initialised *KafkaProducerConfig
func NewKafkaProducerConfig(config *kafka.ConfigMap) *KafkaProducerConfig {
	return &KafkaProducerConfig{
		kafkaConfig: &kafkaConfig{
			config:      config,
			tokenSource: emptyTokenSource{},
			errFn:       func(error) {},
		},
		deliveryErrHandler: func(Event, error) {},
	}
}

// WithEventDeliveryErrHandler registers a delivery error handler to be called
// whenever a delivery fails.
func (pc *KafkaProducerConfig) WithEventDeliveryErrHandler(errHandler func(Event, error)) {
	pc.deliveryErrHandler = errHandler
}

// WithFlushTimeout sets the producer Flush timeout.
func (pc *KafkaProducerConfig) WithFlushTimeout(timeout int) {
	pc.flushTimeoutMs = timeout
}

// NewKafkaProducer returns new a producer.
// To handle errors, either `kafka.Error` messages or any other error while
// interacting with Kafka, register a Error function on *KafkaConsumerConfig.
func NewKafkaProducer(c *KafkaProducerConfig) (Producer, error) {
	kp := &kafkaProducer{
		kafkaCommon: &kafkaCommon{
			errFn:       c.errFn,
			tokenSource: c.tokenSource,
		},
		config: c.config,

		deliveryErrHandler: c.deliveryErrHandler,
		flushTimeoutMs:     c.flushTimeoutMs,
	}

	p, err := kafka.NewProducer(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafkaProducer: %w", err)
	}
	kp.producer = p

	return kp, nil
}

// HandleEvents listens to the producer events channel.
// To register handlers check *KafkaProducerConfig.
func (p *kafkaProducer) HandleEvents() error {
	if p.running() {
		return ErrProducerIsAlreadyRunning
	}
	p.startRunning()

	go func() {
		for kafkaEvent := range p.producer.Events() {
			switch kafkaEvent := kafkaEvent.(type) {
			case *kafka.Message:
				if kafkaEvent.TopicPartition.Error != nil {
					if p.deliveryErrHandler != nil {
						event := messageToEvent(kafkaEvent)
						p.deliveryErrHandler(
							*event,
							fmt.Errorf(
								"%w, topic '%s', partition: '%d', offset: '%s'",
								kafkaEvent.TopicPartition.Error,
								*kafkaEvent.TopicPartition.Topic,
								kafkaEvent.TopicPartition.Partition,
								kafkaEvent.TopicPartition.Offset.String()))
					}
				}
			case kafka.OAuthBearerTokenRefresh:
				p.refreshToken()
			case kafka.Error:
				p.errFn(kafkaEvent)
			}
		}
	}()

	return nil
}

func (p *kafkaProducer) refreshToken() {
	p.kafkaCommon.refreshToken(p.producer)
}

// Sends messages to the given topic. Delivery errors are sent to the producer's
// event channel. To handle delivery errors check *KafkaProducerConfig.WithEventDeliveryErrHandler.
func (p *kafkaProducer) Send(event Event, topic string) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            event.Key,
		Value:          event.Payload,
		Headers:        parseProducerHeaders(event.Headers),
	}, nil)
}

// Shuts the producer down and also closes the underlying kafka producer.
func (p *kafkaProducer) Shutdown(ctx context.Context) error {
	isRunning := p.running()
	p.stopRun()
	defer p.producer.Close()

	// Flush will only work if we listen to producer's events
	if !isRunning {
		return ErrProducerNotHandlingEvents
	}

	for {
		notSent := p.producer.Flush(p.flushTimeoutMs)
		if notSent == 0 {
			break
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%d messages not sent", notSent)
		default:
			continue
		}
	}

	return nil
}

func (p *kafkaProducer) startRunning() bool {
	p.runMu.RLock()
	p.run = true
	defer p.runMu.RUnlock()
	return p.run
}

func (p *kafkaProducer) running() bool {
	p.runMu.RLock()
	defer p.runMu.RUnlock()
	return p.run
}

func (p *kafkaProducer) stopRun() {
	p.runMu.Lock()
	defer p.runMu.Unlock()
	p.run = false
}

func parseProducerHeaders(eventHeaders Header) []kafka.Header {
	var messageHeaders []kafka.Header
	for key, value := range eventHeaders {
		messageHeaders = append(messageHeaders, kafka.Header{Key: key, Value: []byte(value)})
	}
	return messageHeaders
}
