package events

import (
	"context"
	"fmt"
	"log"
	"sync"

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
	activeConn   int
	run          bool
	done         chan struct{}

	handlers []Handler
}

func NewKafkaConsumer(kafkaConsumer *kafka.Consumer, handlers ...Handler) Consumer {
	return KafkaConsumer{
		kafkaConsumer: kafkaConsumer,
		handlers:      handlers,

		activeConnMu: &sync.Mutex{},
		runMu:        &sync.Mutex{},

		activeConn: 0,
		run:        true,
		done:       make(chan struct{}),
	}
}

func (c KafkaConsumer) Run() {
	go func() {
		for c.running() {
			log.Print("waiting message")
			msg, err := c.kafkaConsumer.ReadMessage(-1)
			if err != nil {
				log.Print(fmt.Sprintf("[ERROR] failed to read message: %v", err))
				continue
			}

			for _, h := range c.handlers {
				c.wg.Add(1)
				go func(h Handler) {
					e, err := parseMessage(msg)
					if err != nil {
						// TODO: handle the error!
						return
					}

					// errors are ignored, you handler should
					err = h.Handle(context.Background(), *e)
					c.wg.Done()
				}(h)
			}
		}

		c.wg.Wait()
		close(c.done)
	}()
}

// Shutdown receives a context with deadline or will wait forever
func (c KafkaConsumer) Shutdown(ctx context.Context) error {
	c.stopRun()

	if c.activeConn == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("%d handlers not finished", c.activeConn)
	case <-c.done:
		return nil
	}
}

func (c KafkaConsumer) running() bool {
	c.runMu.Lock()
	defer c.runMu.Unlock()
	return c.run
}

func (c KafkaConsumer) stopRun() {
	c.runMu.Lock()
	defer c.runMu.Unlock()
	c.run = false
}

func parseMessage(m *kafka.Message) (*Event, error) {
	return &Event{}, nil
}
