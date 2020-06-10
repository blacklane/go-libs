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
	Run()
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

func (c *KafkaConsumer) Run() {
	go func() {
		for c.running() {
			msg, err := c.kafkaConsumer.ReadMessage(time.Second)
			if err != nil {
				switch err.(type) {
				case kafka.Error:
					if err.(kafka.Error).Code() != kafka.ErrTimedOut {
						log.Printf("[ERROR] failed to read message: %v", err)
					}
				default:
					log.Printf("[ERROR] failed to read message: %v", err)
				}
				continue
			}

			for _, h := range c.handlers {
				c.wg.Add(1)
				go func(h Handler) {
					defer c.wg.Done()
					e, err := parseMessage(msg)
					if err != nil {
						// TODO: handle the error!
						return
					}

					// errors are ignored, you handler should
					err = h.Handle(context.Background(), *e)
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

func parseMessage(m *kafka.Message) (*Event, error) {
	// TODO: add headers
	return &Event{Payload: m.Value}, nil
}
