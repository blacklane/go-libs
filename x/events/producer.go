package events

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ErrProducerNotHandlingMessages = errors.New("producer should be handling messages")
var ErrProducerIsAlreadyRunning = errors.New("producer is already running")

type Producer interface {
	Send(event Event, topic string) error
	HandleMessages() error
	Shutdown(ctx context.Context) error
}

type ErrorHandler interface {
	HandleError(Event, error)
}

type ErrorHandlerFunc func(Event, error)

func (eh ErrorHandlerFunc) HandleError(e Event, err error) {
	eh(e, err)
}

type kafkaProducer struct {
	runMutex       *sync.Mutex
	producer       *kafka.Producer
	errorHandler   ErrorHandler
	isRunning      bool
	flushTimeoutMs int
}

// NewKafkaProducer returns new a producer.
// It fully manages underlying kafka.Producer's lifecycle.
// Use the 'With...' functions for further configurations:
//   producer, err := NewKafkaProducer(c, errorHandler, WithFlushTimeout(500))
func NewKafkaProducer(
	c *kafka.ConfigMap,
	errorHandler ErrorHandler,
	configs ...func(cfg *kafkaProducer)) (Producer, error) {

	p, err := kafka.NewProducer(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafkaProducer: %w", err)
	}

	kp := &kafkaProducer{
		runMutex:       &sync.Mutex{},
		producer:       p,
		errorHandler:   errorHandler,
		isRunning:      false,
		flushTimeoutMs: 500,
	}

	for _, cfgFn := range configs {
		cfgFn(kp)
	}

	return kp, nil
}

// WithTimeout sets the flush timeout to the given timeout in milliseconds.
func WithFlushTimeout(timeout int) func(p *kafkaProducer) {
	return func(p *kafkaProducer) {
		p.flushTimeoutMs = timeout
	}
}

func (p *kafkaProducer) running() bool {
	p.runMutex.Lock()
	defer p.runMutex.Unlock()
	return p.isRunning
}

func (p *kafkaProducer) stopRunning() {
	p.runMutex.Lock()
	defer p.runMutex.Unlock()
	p.isRunning = false
}

func (p *kafkaProducer) startRunning() {
	p.runMutex.Lock()
	defer p.runMutex.Unlock()
	p.isRunning = true
}

// HandleMessages listens for messages delivery and sends them to error handler
// if the delivery failed.
func (p *kafkaProducer) HandleMessages() error {
	if p.running() {
		return ErrProducerIsAlreadyRunning
	}
	p.startRunning()

	go func() {
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					message := messageToEvent(ev)
					p.errorHandler.HandleError(*message, ev.TopicPartition.Error)
				}
			}
		}
	}()

	return nil
}

// Sends messages. Bear in mind that even if the error is not returned here
// that doesn't mean that the message is delivered (see HandleMessages method).
func (p *kafkaProducer) Send(event Event, topic string) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            event.Key,
		Value:          event.Payload,
	}, nil)
}

// Shuts the producer down and also closes the underlying kafka instance.
func (p *kafkaProducer) Shutdown(ctx context.Context) error {
	isRunning := p.running()
	p.stopRunning()
	defer p.producer.Close()

	// Flush will only work if we listen to producer's events
	if !isRunning {
		return ErrProducerNotHandlingMessages
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
