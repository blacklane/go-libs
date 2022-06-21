package zerolog

import (
	"time"

	"github.com/blacklane/go-libs/logr/field"
	"github.com/rs/zerolog"
)

func contextWithFields(ctx zerolog.Context, fields ...field.Field) zerolog.Context {
	if len(fields) > 0 {
		enc := &contextEncoder{Context: ctx}
		for _, f := range fields {
			f.Encode(enc)
		}
		return enc.Context
	}
	return ctx
}

type contextEncoder struct {
	Context zerolog.Context
}

func (c *contextEncoder) String(key, value string) {
	c.Context = c.Context.Str(key, value)
}

func (c *contextEncoder) Bool(key string, value bool) {
	c.Context = c.Context.Bool(key, value)
}

func (c *contextEncoder) Byte(key string, value byte) {
	c.Context = c.Context.Uint8(key, value)
}

func (c *contextEncoder) Bytes(key string, value []byte) {
	c.Context = c.Context.Bytes(key, value)
}

func (c *contextEncoder) Int8(key string, value int8) {
	c.Context = c.Context.Int8(key, value)
}

func (c *contextEncoder) Int16(key string, value int16) {
	c.Context = c.Context.Int16(key, value)
}

func (c *contextEncoder) Int32(key string, value int32) {
	c.Context = c.Context.Int32(key, value)
}

func (c *contextEncoder) Int64(key string, value int64) {
	c.Context = c.Context.Int64(key, value)
}

func (c *contextEncoder) Uint16(key string, value uint16) {
	c.Context = c.Context.Uint16(key, value)
}

func (c *contextEncoder) Uint32(key string, value uint32) {
	c.Context = c.Context.Uint32(key, value)
}

func (c *contextEncoder) Uint64(key string, value uint64) {
	c.Context = c.Context.Uint64(key, value)
}

func (c *contextEncoder) Float32(key string, value float32) {
	c.Context = c.Context.Float32(key, value)
}

func (c *contextEncoder) Float64(key string, value float64) {
	c.Context = c.Context.Float64(key, value)
}

func (c *contextEncoder) Time(key string, value time.Time) {
	c.Context = c.Context.Time(key, value)
}

func (c *contextEncoder) Duration(key string, value time.Duration) {
	c.Context = c.Context.Dur(key, value)
}

func (c *contextEncoder) Error(value error) {
	c.Context = c.Context.Err(value)
}
