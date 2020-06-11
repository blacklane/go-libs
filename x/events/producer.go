package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

var ErrProducerNotHandlingMessages = errors.New("producer should be handling messages")

const defaultTimeoutMs = 500

type Producer interface {
	Send(event Event, topic string) error
	HandleMessages()
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

func NewKafkaProducer(p *kafka.Producer, errorHandler ErrorHandler) Producer {
	return &kafkaProducer{
		runMutex:       &sync.Mutex{},
		producer:       p,
		errorHandler:   errorHandler,
		isRunning:      false,
		flushTimeoutMs: defaultTimeoutMs,
	}
}

func NewKafkaProducerWithTimeout(p *kafka.Producer, errorHandler ErrorHandler, flushTimeoutMs int) Producer {
	return &kafkaProducer{
		runMutex:       &sync.Mutex{},
		producer:       p,
		errorHandler:   errorHandler,
		isRunning:      false,
		flushTimeoutMs: flushTimeoutMs,
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

func (p *kafkaProducer) HandleMessages() {
	if p.running() {
		return
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
}

func (p *kafkaProducer) Send(event Event, topic string) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            event.Key,
		Value:          event.Payload,
	}, nil)
}

func (p *kafkaProducer) Shutdown(ctx context.Context) error {
	// This method is not responsible for closing the producer
	// as the producer may be reused by different instance
	isRunning := p.running()
	p.stopRunning()

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
