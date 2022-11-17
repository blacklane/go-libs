package consumer

import (
	"encoding/json"
	"time"

	"github.com/blacklane/go-libs/x/events"
)

type Message interface {
	EventName() string
	Header() events.Header
	TopicPartition() events.TopicPartition
	DecodePayload(v interface{}) error
}

type jsonPayload struct {
	Event     string          `json:"event"`
	CreatedAt time.Time       `json:"created_at"`
	Payload   json.RawMessage `json:"payload"`
}

func (m *jsonPayload) DecodePayload(v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}

type jsonMessage struct {
	event   events.Event
	payload jsonPayload
}

func createJsonMessage(ev events.Event) (Message, error) {
	var payload jsonPayload
	if err := json.Unmarshal(ev.Payload, &payload); err != nil {
		return nil, err
	}
	return &jsonMessage{
		event:   ev,
		payload: payload,
	}, nil
}

func (m *jsonMessage) EventName() string                     { return m.payload.Event }
func (m *jsonMessage) Header() events.Header                 { return m.event.Headers }
func (m *jsonMessage) TopicPartition() events.TopicPartition { return m.event.TopicPartition }
func (m *jsonMessage) DecodePayload(v interface{}) error     { return m.payload.DecodePayload(v) }
