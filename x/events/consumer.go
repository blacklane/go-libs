package events

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ErrShutdownTimeout = errors.New("shutdown timeout: not all handlers finished, not closing kafka client")
var ErrConsumerAlreadyShutdown = errors.New("consumer already shutdown")

type Consumer interface {
	Run(timeout time.Duration)
	Shutdown(ctx context.Context) error
}

// KafkaConsumerConfig holds all possible configurations for the kafka consumer.
// Use NewKafkaProducerConfig to initialise it.
// To see the possible configurations, check the its WithXXX methods and
// *kafkaConfig.WithXXX` methods as well
type KafkaConsumerConfig struct {
	*kafkaConfig
}

type kafkaConsumer struct {
	*kafkaCommon

	consumer *kafka.Consumer

	done chan struct{}

	handlers []Handler

	runMu    sync.RWMutex
	run      bool
	shutdown bool
}

// NewKafkaConsumerConfig returns a initialised *KafkaConsumerConfig
func NewKafkaConsumerConfig(config *kafka.ConfigMap) *KafkaConsumerConfig {
	// we do manual commits after processing a message
	autoCommit, err := config.Get("enable.auto.commit", "")
	if err == nil && autoCommit == "" {
		config.SetKey("enable.auto.commit", "false")
	}

	return &KafkaConsumerConfig{
		kafkaConfig: &kafkaConfig{
			config:      config,
			tokenSource: emptyTokenSource{},
			errFn:       func(error) {},
		},
	}
}

// NewKafkaConsumer returns a Consumer which will send every message to all
// handlers and ignore any error returned by them. A middleware should handle
// the errors.
// To handle errors, either `kafka.Error` messages or any other error while
// interacting with Kafka, register a Error function on *KafkaConsumerConfig.
func NewKafkaConsumer(config *KafkaConsumerConfig, topics []string, handlers ...Handler) (Consumer, error) {
	consumer, err := kafka.NewConsumer(config.config)
	if err != nil {
		return nil, fmt.Errorf("could not create kafka consumer: %v", err)
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, fmt.Errorf("could not subscribe to topics: %w", err)
	}

	return &kafkaConsumer{
		kafkaCommon: &kafkaCommon{
			errFn:       config.errFn,
			tokenSource: config.tokenSource,
		},

		consumer: consumer,
		handlers: handlers,

		done: make(chan struct{}),
	}, nil
}

// Run starts to consume messages. If the timeout is -1 it'll block the loop
// until a message arrives. Also during the tests if a timeout = -1 was passed
// the kafka consumer would only read new messages, even if the consumer had
// been created to read from the beginning of the topic.
func (c *kafkaConsumer) Run(timeout time.Duration) {
	c.startRunning()

	go func() {
		timeoutMs := int(timeout.Milliseconds())
		for c.running() {
			kev := c.consumer.Poll(timeoutMs)

			switch kmt := kev.(type) {
			case *kafka.Message:
				c.deliverMessage(kmt)
			case kafka.OAuthBearerTokenRefresh:
				c.refreshToken()
			case kafka.Error:
				if kmt.Code() != kafka.ErrTimedOut {
					c.errFn(fmt.Errorf("failed to read message: %w", kmt))
				}
			case nil: // when c.consumer.Poll(timeoutMs) times out, it returns nil.
				continue
			default:
				c.errFn(fmt.Errorf("unknown kafka message type: '%T': %#v", kev, kev))
			}
		}

		close(c.done)
	}()
}

// Shutdown receives a context with deadline, or it will wait until all handlers to
// finish. The closing of the underlying kafka client is skipped if the cxt is
// done before trying to close it. Also, the underlying kafka client close does
// not respect the ctx cancellation.
// If Shutdown is called more than once, it immediately returns ErrConsumerAlreadyShutdown
// for the subsequent calls.
func (c *kafkaConsumer) Shutdown(ctx context.Context) error {
	if c.shutdown {
		return ErrConsumerAlreadyShutdown
	}
	c.shutdown = true

	c.stopRun()
	var err error
	select {
	case <-ctx.Done():
		return ErrShutdownTimeout
	case <-c.done:
		err = nil
	}

	errClose := c.consumer.Close()
	if errClose != nil {
		errClose = fmt.Errorf("close kafka consumer failed: %w", errClose)
	}

	if err != nil {
		return fmt.Errorf("shutdown failures: %w, %s", err, errClose)
	}

	return nil
}

func (c *kafkaConsumer) deliverMessage(msg *kafka.Message) {
	e := messageToEvent(msg)

	waitGroup := sync.WaitGroup{}

	worker := func(handler Handler) {
		defer waitGroup.Done()

		err := handler.Handle(context.Background(), *e)
		if err != nil {
			c.errFn(err)
		}
	}

	for _, h := range c.handlers {
		waitGroup.Add(1)
		go worker(h)
	}

	waitGroup.Wait()

	_, err := c.consumer.Commit()
	if err != nil {
		c.errFn(fmt.Errorf("error commiting offset: %v", err))
	}
}

func (c *kafkaConsumer) refreshToken() {
	c.kafkaCommon.refreshToken(c.consumer)
}

func messageToEvent(m *kafka.Message) *Event {
	return &Event{Payload: m.Value, Headers: parseHeaders(m.Headers), Key: m.Key}
}

func parseHeaders(headers []kafka.Header) Header {
	hs := Header{}
	for _, kh := range headers {
		hs[kh.Key] = string(kh.Value)
	}

	return hs
}

func (c *kafkaConsumer) startRunning() bool {
	c.runMu.RLock()
	defer c.runMu.RUnlock()
	c.run = true
	return c.run
}

func (c *kafkaConsumer) running() bool {
	c.runMu.RLock()
	defer c.runMu.RUnlock()
	return c.run
}

func (c *kafkaConsumer) stopRun() {
	c.runMu.Lock()
	defer c.runMu.Unlock()
	c.run = false
}
