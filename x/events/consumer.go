package events

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer interface {
	Run(timeout time.Duration)
	Shutdown(ctx context.Context) error
}

type KafkaConsumer struct {
	kafkaConsumer *kafka.Consumer

	wg           *sync.WaitGroup
	activeConnMu *sync.Mutex
	runMu        *sync.Mutex
	run          bool
	done         chan struct{}

	handlers []Handler
}

func NewKafkaConsumer(kafkaConsumer *kafka.Consumer, handlers ...Handler) Consumer {
	return &KafkaConsumer{
		kafkaConsumer: kafkaConsumer,
		handlers:      handlers,

		wg:           &sync.WaitGroup{},
		activeConnMu: &sync.Mutex{},
		runMu:        &sync.Mutex{},

		run:  true,
		done: make(chan struct{}),
	}
}

func (c *KafkaConsumer) Run(timeout time.Duration) {
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
					log.Printf("[ERROR] failed to read message: %v", err)
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

// Shutdown receives a context with deadline or will wait forever
func (c *KafkaConsumer) Shutdown(ctx context.Context) error {
	c.stopRun()

	select {
	case <-ctx.Done():
		return errors.New("not all handlers finished")
	case <-c.done:
		return nil
	}
}

func (c KafkaConsumer) running() bool {
	c.runMu.Lock()
	defer c.runMu.Unlock()
	return c.run
}

func (c *KafkaConsumer) stopRun() {
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
