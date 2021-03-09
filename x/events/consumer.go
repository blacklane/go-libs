package events

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ErrShutdownTimeout = errors.New("shutdown timeout: not all handlers finished")
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
	orderKey OrderKey
}

type OrderKey int

const (
	OrderByMessageKey OrderKey = iota
	OrderByTrackingId
	OrderByTimestamp
	OrderByOffset
	OrderByNotSpecified
)

type kafkaConsumer struct {
	*kafkaCommon

	consumer *kafka.Consumer

	wg           *sync.WaitGroup
	activeConnMu sync.Mutex
	done         chan struct{}

	handlers []Handler

	runMu    sync.RWMutex
	run      bool
	shutdown bool

	orderKey       OrderKey
	keysInProgress sync.Map
}

// NewKafkaConsumerConfig returns a initialised *KafkaConsumerConfig
func NewKafkaConsumerConfig(config *kafka.ConfigMap) *KafkaConsumerConfig {
	return &KafkaConsumerConfig{
		kafkaConfig: &kafkaConfig{
			config:      config,
			tokenSource: emptyTokenSource{},
			errFn:       func(error) {},
		},
		orderKey: OrderByNotSpecified,
	}
}

func (k *KafkaConsumerConfig) WithOrder(orderKey OrderKey) *KafkaConsumerConfig {
	k.orderKey = orderKey
	return k
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

		wg:           &sync.WaitGroup{},
		activeConnMu: sync.Mutex{},

		done:     make(chan struct{}),
		orderKey: config.orderKey,
	}, nil
}

// Run starts to consume messages. If the timeout is -1 it'll block the loop
// until a message arrives. Also during the tests if a timeout = -1 was passed
// the kafka consumer would only read new messages, even if the consumer had
// been created to read from the beginning of the topic.
func (c *kafkaConsumer) Run(timeout time.Duration) {
	c.startRunning()
	go func() {
		for c.running() {
			kev := c.consumer.Poll(int(timeout.Milliseconds()))
			switch kev := kev.(type) {
			case *kafka.Message:
				c.deliverMessage(kev)
			case kafka.OAuthBearerTokenRefresh:
				c.refreshToken()
			case kafka.Error:
				if kev.Code() != kafka.ErrTimedOut {
					c.errFn(fmt.Errorf("failed to read message: %w", kev))
				}
			}
		}

		c.wg.Wait()
		close(c.done)
	}()
}

// Shutdown receives a context with deadline or will wait until all handlers to
// finish. Shutdown can only be called once, if called more than once it returns
// ErrConsumerAlreadyShutdown
func (c *kafkaConsumer) Shutdown(ctx context.Context) error {
	if c.shutdown {
		return ErrConsumerAlreadyShutdown
	}

	c.stopRun()

	var err error
	select {
	case <-ctx.Done():
		err = ErrShutdownTimeout
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
	c.wg.Add(1)
	e := messageToEvent(msg)

	var orderKey string
	var orderMutex sync.Locker
	switch c.orderKey {
	case OrderByNotSpecified:
		orderMutex = &sync.Mutex{} // always new mutex
	case OrderByMessageKey:
		orderKey = string(e.Key)
	case OrderByTrackingId:
		orderKey = e.Headers[HeaderTrackingID]
	}

	if orderMutex == nil {
		keyMutex, _ := c.keysInProgress.LoadOrStore(orderKey, &sync.Mutex{})
		orderMutex = keyMutex.(sync.Locker)
	}

	go func(handlers []Handler, orderKey string, keysInProgress *sync.Map, orderMutex sync.Locker) {
		defer c.wg.Done()
		orderMutex.Lock()

		// Errors are ignored, a middleware or the handler should handle them
		for _, h := range handlers {
			_ = h.Handle(context.Background(), *e)
		}
		orderMutex.Unlock()
		keysInProgress.Delete(orderKey)
	}(c.handlers, orderKey, &c.keysInProgress, orderMutex)
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
	c.run = true
	defer c.runMu.RUnlock()
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
