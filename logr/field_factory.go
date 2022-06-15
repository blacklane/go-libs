package logr

import "time"

func String(key string, value string) Field {
	return &field{
		Key:    key,
		String: &value,
		Kind:   kindString,
	}
}

func Bool(key string, value bool) Field {
	return &field{
		Key:  key,
		Bool: &value,
		Kind: kindBool,
	}
}

func Byte(key string, value byte) Field {
	unsigned := uint64(value)
	return &field{
		Key:      key,
		Unsigned: &unsigned,
		Kind:     kindByte,
	}
}

func Bytes(key string, value []byte) Field {
	return &field{
		Key:   key,
		Value: value,
		Kind:  kindBytes,
	}
}

func Int8(key string, value int8) Field {
	signed := int64(value)
	return &field{
		Key:    key,
		Signed: &signed,
		Kind:   kindInt8,
	}
}

func Int16(key string, value int16) Field {
	signed := int64(value)
	return &field{
		Key:    key,
		Signed: &signed,
		Kind:   kindInt16,
	}
}

func Int32(key string, value int32) Field {
	signed := int64(value)
	return &field{
		Key:    key,
		Signed: &signed,
		Kind:   kindInt32,
	}
}

func Int64(key string, value int64) Field {
	return &field{
		Key:    key,
		Signed: &value,
		Kind:   kindInt64,
	}
}

func Uint16(key string, value uint16) Field {
	unsigned := uint64(value)
	return &field{
		Key:      key,
		Unsigned: &unsigned,
		Kind:     kindUint16,
	}
}

func Uint32(key string, value uint32) Field {
	unsigned := uint64(value)
	return &field{
		Key:      key,
		Unsigned: &unsigned,
		Kind:     kindUint32,
	}
}

func Uint64(key string, value uint64) Field {
	return &field{
		Key:      key,
		Unsigned: &value,
		Kind:     kindUint64,
	}
}

func Float32(key string, value float32) Field {
	v := float64(value)
	return &field{
		Key:   key,
		Float: &v,
		Kind:  kindFloat32,
	}
}

func Float64(key string, value float64) Field {
	return &field{
		Key:   key,
		Float: &value,
		Kind:  kindFloat64,
	}
}

func Time(key string, value time.Time) Field {
	return &field{
		Key:  key,
		Time: &value,
		Kind: kindTime,
	}
}

func Duration(key string, value time.Duration) Field {
	signed := int64(value)
	return &field{
		Key:    key,
		Signed: &signed,
		Kind:   kindDuration,
	}
}

func Error(err error) Field {
	return &field{
		Value: err,
		Kind:  kindError,
	}
}
