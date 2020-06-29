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
	HandleEvents() error
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
	runMutex *sync.Mutex
	producer *kafka.Producer

	errorHandler      ErrorHandler
	messageHandler    func(*kafka.Message)
	kafkaEventHandler func(event kafka.Event)

	isRunning      bool
	flushTimeoutMs int
}

// NewKafkaProducer returns new a producer.
// It fully manages underlying kafka.Producer's lifecycle.
// Use the 'With...' functions for further configurations:
//   producer, err := NewKafkaProducer(c, errorHandler, WithFlushTimeout(500))
func NewKafkaProducer(c *kafka.ConfigMap, configs ...func(*kafkaProducer)) (Producer, error) {

	p, err := kafka.NewProducer(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafkaProducer: %w", err)
	}

	kp := &kafkaProducer{
		runMutex:       &sync.Mutex{},
		producer:       p,
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

// WithErrorHandler adds a handler to handle *kafka.Message when
// Message.TopicPartition.Error != nil
// See also WithKafkaEventHandler and WithKafkaMessageHandler
func WithErrorHandler(errHandler ErrorHandler) func(p *kafkaProducer) {
	return func(p *kafkaProducer) {
		p.errorHandler = errHandler
	}
}

// WithErrorHandler adds a handler to handle *kafka.Message if
// Message.TopicPartition.Error == nil
// See also WithKafkaEventHandler and WithKafkaMessageHandler
func WithKafkaMessageHandler(errHandler ErrorHandler) func(p *kafkaProducer) {
	return func(p *kafkaProducer) {
		p.errorHandler = errHandler
	}
}

//
func WithKafkaEventHandler(errHandler ErrorHandler) func(p *kafkaProducer) {
	return func(p *kafkaProducer) {
		p.errorHandler = errHandler
	}
}

// HandleEvents listens for messages delivery and sends them to error handler
// if the delivery failed.
func (p *kafkaProducer) HandleEvents() error {
	if p.running() {
		return ErrProducerIsAlreadyRunning
	}
	p.startRunning()

	go func() {
		for kafkaEvent := range p.producer.Events() {
			switch kafkaEvent := kafkaEvent.(type) {
			case *kafka.Message:
				p.handleMessage(kafkaEvent, kafkaEvent)
			default:
				if p.kafkaEventHandler != nil {
					p.kafkaEventHandler(kafkaEvent)
				}
			}
		}
	}()

	return nil
}

// handleMessage calls the appropriated handler, if there is no handler for
// *kafka.Message, the generic kafkaEventHandler is invoked
func (p *kafkaProducer) handleMessage(e kafka.Event, msg *kafka.Message) {
	if msg.TopicPartition.Error != nil {
		event := messageToEvent(msg)
		p.errorHandler.HandleError(
			*event,
			msg.TopicPartition.Error)
		return
	}

	if p.messageHandler != nil {
		p.messageHandler(msg)
		return
	}

	if p.kafkaEventHandler != nil {
		p.kafkaEventHandler(e)
		return
	}

	return
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
