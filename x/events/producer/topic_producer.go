package producer

import (
	"context"
	"fmt"

	"github.com/blacklane/go-libs/x/events"
	"go.opentelemetry.io/otel"
)

type (
	EventHook     func(*events.Event)
	TopicProducer interface {
		Send(ctx context.Context, eventName string, payload interface{}, hooks ...EventHook) error
	}
)

func WithKey(key string) EventHook {
	return func(e *events.Event) {
		e.Key = []byte(key)
	}
}

type topicProducer struct {
	producer events.Producer
	topic    string
}

func NewTopicProducer(producer events.Producer, topic string) TopicProducer {
	return &topicProducer{
		producer: producer,
		topic:    topic,
	}
}

func (p *topicProducer) Send(ctx context.Context, eventName string, payload interface{}, hooks ...EventHook) error {
	bytes, err := crateEventPayload(eventName, payload)
	if err != nil {
		return fmt.Errorf("could not marshal kafka event payload: %w", err)
	}

	ev := events.Event{
		Headers: events.Header{},
		Payload: bytes,
	}

	otel.GetTextMapPropagator().Inject(ctx, ev.Headers)

	for _, hook := range hooks {
		hook(&ev)
	}

	return p.producer.SendCtx(ctx, eventName, ev, p.topic)
}
