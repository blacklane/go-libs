package consumer

import (
	"fmt"
	"testing"

	"github.com/blacklane/go-libs/x/events"
)

func TestCreateJsonMessage(t *testing.T) {
	eventName := "event-name"
	payloadName := "payload-name"

	ev := events.Event{
		Payload: []byte(fmt.Sprintf(`{
			"event": "%s",
			"payload": {
				"name": "%s"
			}
		}`, eventName, payloadName)),
	}

	m, err := createJsonMessage(ev)
	if err != nil {
		t.Fatalf("error creating message: %v", err)
	}

	if m.EventName() != eventName {
		t.Errorf("invalid name, want: %s, got: %s", eventName, m.EventName())
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := m.DecodePayload(&payload); err != nil {
		t.Fatalf("error decoding payload: %v", err)
	}

	if payload.Name != payloadName {
		t.Errorf("invalid payload name, want: %s, got: %s", payloadName, payload.Name)
	}
}
