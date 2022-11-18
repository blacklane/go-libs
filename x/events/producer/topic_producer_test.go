package producer

import (
	"context"
	"testing"

	"github.com/blacklane/go-libs/x/events/eventstest"
	"github.com/golang/mock/gomock"
)

type eventPayload struct {
	Value string
}

func TestTopicProducer_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	p := eventstest.NewMockProducer(ctrl)

	topic := "topic-name"
	eventName := "event-name"
	eventKey := "event-key"
	ctx := context.Background()
	payload := eventPayload{
		Value: "1234",
	}

	p.EXPECT().SendCtx(
		ctx,
		eventName,
		eventstest.MatchEvent([]byte(eventKey), eventName, payload),
		topic,
	).Return(nil).Times(1)

	tp := NewTopicProducer(p, topic)

	if err := tp.Send(ctx, eventName, payload, WithKey(eventKey)); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
