package field

import (
	"fmt"
	"time"
)

type Field interface {
	Encode(Encoder)
}

type kind uint8

const (
	kindString kind = iota
	kindBool
	kindByte
	kindBytes
	kindInt8
	kindInt16
	kindInt32
	kindInt64
	kindUint8
	kindUint16
	kindUint32
	kindUint64
	kindFloat32
	kindFloat64
	kindTime
	kindDuration
	kindError
)

type field struct {
	Key      string
	Bool     *bool
	String   *string
	Signed   *int64
	Unsigned *uint64
	Float    *float64
	Time     *time.Time
	Value    interface{}
	Kind     kind
}

func (f *field) Encode(enc Encoder) {
	switch f.Kind {
	case kindString:
		enc.String(f.Key, *f.String)
	case kindBool:
		enc.Bool(f.Key, *f.Bool)
	case kindByte:
		enc.Byte(f.Key, byte(*f.Unsigned))
	case kindBytes:
		enc.Bytes(f.Key, f.Value.([]byte))
	case kindInt8:
		enc.Int8(f.Key, int8(*f.Signed))
	case kindInt16:
		enc.Int16(f.Key, int16(*f.Signed))
	case kindInt32:
		enc.Int32(f.Key, int32(*f.Signed))
	case kindInt64:
		enc.Int64(f.Key, *f.Signed)
	case kindUint16:
		enc.Uint16(f.Key, uint16(*f.Unsigned))
	case kindUint32:
		enc.Uint32(f.Key, uint32(*f.Unsigned))
	case kindUint64:
		enc.Uint64(f.Key, uint64(*f.Unsigned))
	case kindFloat32:
		enc.Float32(f.Key, float32(*f.Float))
	case kindFloat64:
		enc.Float64(f.Key, *f.Float)
	case kindTime:
		enc.Time(f.Key, *f.Time)
	case kindDuration:
		enc.Duration(f.Key, time.Duration(*f.Signed))
	case kindError:
		enc.Error(f.Value.(error))
	default:
		panic(fmt.Sprintf("unknown field type: %d", f.Kind))
	}
}
