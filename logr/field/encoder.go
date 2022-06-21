package field

import "time"

type Encoder interface {
	String(key, value string)
	Bool(key string, value bool)
	Byte(key string, value byte)
	Bytes(key string, value []byte)
	Int8(key string, value int8)
	Int16(key string, value int16)
	Int32(key string, value int32)
	Int64(key string, value int64)
	Uint16(key string, value uint16)
	Uint32(key string, value uint32)
	Uint64(key string, value uint64)
	Float32(key string, value float32)
	Float64(key string, value float64)
	Time(key string, value time.Time)
	Duration(key string, value time.Duration)
	Error(err error)
}
