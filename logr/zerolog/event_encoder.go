package zerolog

import (
	"time"

	"github.com/blacklane/go-libs/logr/field"
	"github.com/rs/zerolog"
)

func eventWithFields(ev *zerolog.Event, fields ...field.Field) *zerolog.Event {
	if len(fields) > 0 {
		enc := &eventEncoder{Event: ev}
		for _, f := range fields {
			f.Encode(enc)
		}
		return enc.Event
	}
	return ev
}

type eventEncoder struct {
	Event *zerolog.Event
}

func (e *eventEncoder) String(key, value string) {
	e.Event = e.Event.Str(key, value)
}

func (e *eventEncoder) Bool(key string, value bool) {
	e.Event = e.Event.Bool(key, value)
}

func (e *eventEncoder) Byte(key string, value byte) {
	e.Event = e.Event.Uint8(key, value)
}

func (e *eventEncoder) Bytes(key string, value []byte) {
	e.Event = e.Event.Bytes(key, value)
}

func (e *eventEncoder) Int8(key string, value int8) {
	e.Event = e.Event.Int8(key, value)
}

func (e *eventEncoder) Int16(key string, value int16) {
	e.Event = e.Event.Int16(key, value)
}

func (e *eventEncoder) Int32(key string, value int32) {
	e.Event = e.Event.Int32(key, value)
}

func (e *eventEncoder) Int64(key string, value int64) {
	e.Event = e.Event.Int64(key, value)
}

func (e *eventEncoder) Uint16(key string, value uint16) {
	e.Event = e.Event.Uint16(key, value)
}

func (e *eventEncoder) Uint32(key string, value uint32) {
	e.Event = e.Event.Uint32(key, value)
}

func (e *eventEncoder) Uint64(key string, value uint64) {
	e.Event = e.Event.Uint64(key, value)
}

func (e *eventEncoder) Float32(key string, value float32) {
	e.Event = e.Event.Float32(key, value)
}

func (e *eventEncoder) Float64(key string, value float64) {
	e.Event = e.Event.Float64(key, value)
}

func (e *eventEncoder) Time(key string, value time.Time) {
	e.Event = e.Event.Time(key, value)
}

func (e *eventEncoder) Duration(key string, value time.Duration) {
	e.Event = e.Event.Dur(key, value)
}

func (e *eventEncoder) Error(value error) {
	e.Event = e.Event.Err(value)
}
