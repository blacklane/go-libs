package events

import (
	"context"
	"errors"
	"fmt"
	"log"
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

type kafkaConsumer struct {
	kafkaConsumer *kafka.Consumer

	wg           *sync.WaitGroup
	activeConnMu sync.Mutex
	runMu        sync.RWMutex
	run          bool
	shutdown     bool
	done         chan struct{}

	handlers []Handler
}

// NewKafkaConsumer returns a Consumer which will send every message to all
// handlers and ignore any error returned by them. A middleware should handle
// the errors.
// If the kafka consumer receives a kafka.Error it'll log it using the log
// package.
// TODO: improve error handling. Ideas:
//  - Receive an io.Writer
//  - go-libs/logger.Logger
//  - create an errors channel to send all errors launching goroutines to send
//  the errors so Run won't block if the channel is not drained
func NewKafkaConsumer(config *kafka.ConfigMap, topics []string, handlers ...Handler) (Consumer, error) {
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("could not create kafka consumer: %v", err)
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, fmt.Errorf("could not subscribe to topics: %w", err)
	}

	return &kafkaConsumer{
		kafkaConsumer: consumer,
		handlers:      handlers,

		wg:           &sync.WaitGroup{},
		activeConnMu: sync.Mutex{},
		runMu:        sync.RWMutex{},

		run:  true,
		done: make(chan struct{}),
	}, nil
}

// Run starts to consume messages. If the timeout is -1 it'll block the loop
// until a message arrives. Also during the tests if a timeout = -1 was passed
// the kafka consumer would only read new messages, even if the consumer had
// been created to read from the beginning of the topic.
func (c *kafkaConsumer) Run(timeout time.Duration) {
	go func() {
		for c.running() {
			msg, err := c.kafkaConsumer.ReadMessage(timeout)
			if err != nil {
				switch err.(type) {
				case kafka.Error:
					if err.(kafka.Error).Code() != kafka.ErrTimedOut {
						// TODO: handle it properly!!
						log.Printf("[ERROR] failed to read message: %v", err)
					}
				default:
					// TODO: handle it properly!!
					if msg != nil {
						log.Printf("[ERROR] failed to read message: %v, err: %v", msg, err)
					} else {
						log.Printf("[ERROR] failed to read message: %v", err)
					}
				}
				continue
			}

			for _, h := range c.handlers {
				c.wg.Add(1)
				go func(h Handler) {
					defer c.wg.Done()
					e := messageToEvent(msg)

					// Errors are ignored, a middleware should handle them
					_ = h.Handle(context.Background(), *e)
				}(h)
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

	errClose := c.kafkaConsumer.Close()
	if errClose != nil {
		errClose = fmt.Errorf("close kafka consumer failed: %w", errClose)
	}

	if err != nil {
		return fmt.Errorf("shutdown failures: %w, %s", err, errClose)
	}

	return nil
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

func messageToEvent(m *kafka.Message) *Event {
	return &Event{Payload: m.Value, Headers: parseHeaders(m.Headers)}
}

func parseHeaders(headers []kafka.Header) Header {
	hs := Header{}
	for _, kh := range headers {
		hs[kh.Key] = string(kh.Value)
	}

	return hs
}
