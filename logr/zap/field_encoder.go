package zap

import (
	"time"

	"github.com/blacklane/go-libs/logr/field"
	"go.uber.org/zap"
)

func mapFields(fields []field.Field) []zap.Field {
	enc := &fieldEncoder{}
	for _, f := range fields {
		f.Encode(enc)
	}
	return enc.Fields
}

type fieldEncoder struct {
	Fields []zap.Field
}

func (f *fieldEncoder) append(field zap.Field) {
	f.Fields = append(f.Fields, field)
}

func (f *fieldEncoder) String(key, value string) {
	f.append(zap.String(key, value))
}

func (f *fieldEncoder) Bool(key string, value bool) {
	f.append(zap.Bool(key, value))
}

func (f *fieldEncoder) Byte(key string, value byte) {
	f.append(zap.Uint8(key, value))
}

func (f *fieldEncoder) Bytes(key string, value []byte) {
	f.append(zap.Uint8s(key, value))
}

func (f *fieldEncoder) Int8(key string, value int8) {
	f.append(zap.Int8(key, value))
}

func (f *fieldEncoder) Int16(key string, value int16) {
	f.append(zap.Int16(key, value))
}

func (f *fieldEncoder) Int32(key string, value int32) {
	f.append(zap.Int32(key, value))
}

func (f *fieldEncoder) Int64(key string, value int64) {
	f.append(zap.Int64(key, value))
}

func (f *fieldEncoder) Uint16(key string, value uint16) {
	f.append(zap.Uint16(key, value))
}

func (f *fieldEncoder) Uint32(key string, value uint32) {
	f.append(zap.Uint32(key, value))
}

func (f *fieldEncoder) Uint64(key string, value uint64) {
	f.append(zap.Uint64(key, value))
}

func (f *fieldEncoder) Float32(key string, value float32) {
	f.append(zap.Float32(key, value))
}

func (f *fieldEncoder) Float64(key string, value float64) {
	f.append(zap.Float64(key, value))
}

func (f *fieldEncoder) Time(key string, value time.Time) {
	f.append(zap.Time(key, value))
}

func (f *fieldEncoder) Duration(key string, value time.Duration) {
	f.append(zap.Duration(key, value))
}

func (f *fieldEncoder) Error(value error) {
	f.append(zap.Error(value))
}
