package events

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

type Producer interface {
	Send(event Event) error
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
	topic        string
	run          bool
}

func NewKafkaProducer(p *kafka.Producer, topic string, errorHandler ErrorHandler) Producer {
	return &KafkaProducer{
		runMu:        &sync.Mutex{},
		producer:     p,
		errorHandler: errorHandler,
		topic:        topic,
		run:          false,
	}
}

func (p KafkaProducer) running() bool {
	p.runMu.Lock()
	defer p.runMu.Unlock()
	return p.run
}

func (p KafkaProducer) stopRun() {
	p.runMu.Lock()
	defer p.runMu.Unlock()
	p.run = false
}

func (p KafkaProducer) Run() {
	if p.running() {
		return
	}
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

func (p KafkaProducer) Send(event Event) error {
	if !p.running() {
		return fmt.Errorf("producer should be running")
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            event.Key,
		Value:          event.Payload,
	}, nil)
}

func (p KafkaProducer) Shutdown(ctx context.Context) error {
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
