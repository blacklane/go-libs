package consumer_test

import (
	"context"
	"testing"

	"github.com/blacklane/go-libs/x/events"
	"github.com/blacklane/go-libs/x/events/consumer"
)

type eventPayload struct {
	ID string `json:"id"`
}

func TestConsumer(t *testing.T) {
	tests := []struct {
		event       string
		jsonPayload string
		payloadID   string
		skip        bool
		err         bool
	}{
		{
			event: "skip-event",
			jsonPayload: `{
				"event": "event-name"
			}`,
			skip: true,
		},
		{
			event: "invalid-payload",
			jsonPayload: `{
				"event": "invalid-payload",
				"payload": "invalid-type"
			}`,
			err: true,
		},
		{
			event:     "valid-event",
			payloadID: "1234",
			jsonPayload: `{
				"event": "valid-event",
				"payload": {
					"id": "1234"
				}
			}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.event, func(t *testing.T) {
			t.Parallel()

			payloadID := ""
			handler := consumer.New(test.event, func(ctx context.Context, m consumer.Message) error {
				if test.skip {
					t.Errorf("the event should be skiped: %s", m.EventName())
				}

				var payload eventPayload
				if err := m.DecodePayload(&payload); err != nil {
					return err
				}
				payloadID = payload.ID
				return nil
			})

			err := handler.Handle(context.Background(), events.Event{
				Payload: []byte(test.jsonPayload),
			})

			if test.err && err == nil {
				t.Error("error expected")
			}

			if test.payloadID != payloadID {
				t.Errorf("invalid payload id, want: %s, got: %s", test.payloadID, payloadID)
			}
		})
	}
}
