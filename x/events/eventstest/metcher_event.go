package eventstest

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/blacklane/go-libs/x/events"
	"github.com/golang/mock/gomock"
)

type eventMatcher struct {
	Key          []byte
	EventName    string
	EventPayload interface{}
}

func MatchEvent(key []byte, eventName string, payload interface{}) gomock.Matcher {
	return &eventMatcher{
		Key:          key,
		EventName:    eventName,
		EventPayload: payload,
	}
}

func (m *eventMatcher) Matches(x interface{}) bool {
	ev, ok := x.(events.Event)
	if !ok {
		return false
	}
	if !bytes.Equal(ev.Key, m.Key) {
		return false
	}
	var msg struct {
		Event   string          `json:"event"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(ev.Payload, &msg); err != nil {
		return false
	}
	if msg.Event != m.EventName {
		return false
	}

	if data, err := json.Marshal(m.EventPayload); err != nil || !bytes.Equal(data, msg.Payload) {
		return false
	}

	return true
}

func (m *eventMatcher) String() string {
	return fmt.Sprintf("key: %s, event: %s, payload: %v", m.Key, m.EventName, m.EventPayload)
}
