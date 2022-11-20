package producer

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type jsonMessage struct {
	Event     string      `json:"event"`
	UUID      string      `json:"uuid"`
	CreatedAt time.Time   `json:"created_at"`
	Payload   interface{} `json:"payload"`
}

func crateEventPayload(eventName string, payload interface{}) ([]byte, error) {
	m := jsonMessage{
		Event:     eventName,
		UUID:      uuid.New().String(),
		CreatedAt: time.Now().UTC().Round(time.Second),
		Payload:   payload,
	}
	return json.Marshal(m)
}
