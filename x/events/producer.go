package events

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

type Producer interface {
	SendEvent(event Event) error
	Run()
	Close(ctx context.Context) error
}

type EventHandler func(*Event, error)

type KafkaProducer struct {
	runMu           *sync.Mutex
	producer        *kafka.Producer
	deliveryHandler EventHandler
	topic           string
	run             bool
}

func NewKafkaProducer(p *kafka.Producer, topic string, deliveryHandler EventHandler) Producer {
	return &KafkaProducer{
		runMu:           &sync.Mutex{},
		producer:        p,
		deliveryHandler: deliveryHandler,
		topic:           topic,
		run:             false,
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
				if p.deliveryHandler != nil {
					message, _ := parseMessage(ev)
					p.deliveryHandler(message, ev.TopicPartition.Error)
				}
			}
		}
	}()
}

func (p KafkaProducer) SendEvent(event Event) error {
	if !p.running() {
		return fmt.Errorf("producer should be running")
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            event.Name,
		Value:          event.Payload,
	}, nil)
}

func (p KafkaProducer) Close(ctx context.Context) error {
	if !p.running() {
		return fmt.Errorf("producer should be running")
	}
	p.stopRun()
	var err error
	for {
		notSent := p.producer.Flush(1000)
		if notSent == 0 {
			break
		}
		select {
		case <-ctx.Done():
			err = fmt.Errorf("%d messages not sent", notSent)
		default:
			continue
		}
	}

	p.producer.Close()

	return err
}
