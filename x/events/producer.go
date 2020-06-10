package events

import (
	"context"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer interface {
	Send(event Event, topic string) error
	Run()
	Shutdown(ctx context.Context) error
}

type ErrorHandler interface {
	HandleError(*Event, error)
}

type ErrorHandlerFunc func(*Event, error)

func (eh ErrorHandlerFunc) HandleError(e *Event, err error) {
	eh(e, err)
}

type KafkaProducer struct {
	runMu        *sync.Mutex
	producer     *kafka.Producer
	errorHandler ErrorHandler
	run          bool
}

func NewKafkaProducer(p *kafka.Producer, errorHandler ErrorHandler) Producer {
	return &KafkaProducer{
		runMu:        &sync.Mutex{},
		producer:     p,
		errorHandler: errorHandler,
		run:          false,
	}
}

func (p *KafkaProducer) running() bool {
	p.runMu.Lock()
	defer p.runMu.Unlock()
	return p.run
}

func (p *KafkaProducer) stopRun() {
	p.runMu.Lock()
	defer p.runMu.Unlock()
	p.run = false
}

func (p *KafkaProducer) startRun() {
	p.runMu.Lock()
	defer p.runMu.Unlock()
	p.run = true
}

func (p *KafkaProducer) Run() {
	if p.running() {
		return
	}
	p.startRun()
	go func() {
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if p.errorHandler != nil && ev.TopicPartition.Error != nil {
					message, _ := parseMessage(ev)
					p.errorHandler.HandleError(message, ev.TopicPartition.Error)
				}
			}
		}
	}()
}

func (p *KafkaProducer) Send(event Event, topic string) error {
	if !p.running() {
		return fmt.Errorf("producer should be running")
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            event.Key,
		Value:          event.Payload,
	}, nil)
}

func (p *KafkaProducer) Shutdown(ctx context.Context) error {
	if !p.running() {
		return fmt.Errorf("producer should be running")
	}
	defer p.producer.Close()
	p.stopRun()
	for {
		notSent := p.producer.Flush(1000)
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
