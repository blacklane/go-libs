package logr

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// - encoder mock -

type encoderMock struct {
	mock.Mock
}

func (m *encoderMock) String(key, value string) {
	m.Called(key, value)
}

func (m *encoderMock) Bool(key string, value bool) {
	m.Called(key, value)
}

func (m *encoderMock) Byte(key string, value byte) {
	m.Called(key, value)
}

func (m *encoderMock) Bytes(key string, value []byte) {
	m.Called(key, value)
}

func (m *encoderMock) Int8(key string, value int8) {
	m.Called(key, value)
}

func (m *encoderMock) Int16(key string, value int16) {
	m.Called(key, value)
}

func (m *encoderMock) Int32(key string, value int32) {
	m.Called(key, value)
}

func (m *encoderMock) Int64(key string, value int64) {
	m.Called(key, value)
}

func (m *encoderMock) Uint16(key string, value uint16) {
	m.Called(key, value)
}

func (m *encoderMock) Uint32(key string, value uint32) {
	m.Called(key, value)
}

func (m *encoderMock) Uint64(key string, value uint64) {
	m.Called(key, value)
}

func (m *encoderMock) Float32(key string, value float32) {
	m.Called(key, value)
}

func (m *encoderMock) Float64(key string, value float64) {
	m.Called(key, value)
}

func (m *encoderMock) Time(key string, value time.Time) {
	m.Called(key, value)
}

func (m *encoderMock) Duration(key string, value time.Duration) {
	m.Called(key, value)
}

func (m *encoderMock) Error(err error) {
	m.Called(err)
}

// - factory tests -

func FuzzString(f *testing.F) {
	f.Add("key1", "value1")
	f.Add("key2", "value2")
	f.Add("key3", "value3")

	f.Fuzz(func(t *testing.T, key string, value string) {
		field := String(key, value)

		enc := &encoderMock{}
		enc.On("String", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzBool(f *testing.F) {
	f.Add("key1", true)
	f.Add("key2", false)

	f.Fuzz(func(t *testing.T, key string, value bool) {
		field := Bool(key, value)

		enc := &encoderMock{}
		enc.On("Bool", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzByte(f *testing.F) {
	f.Add("key1", byte(16))
	f.Add("key2", byte(23))
	f.Add("key3", byte(63))

	f.Fuzz(func(t *testing.T, key string, value byte) {
		field := Byte(key, value)

		enc := &encoderMock{}
		enc.On("Byte", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzBytes(f *testing.F) {
	f.Add("key1", []byte{1, 2, 3})
	f.Add("key2", []byte{4, 5, 6})
	f.Add("key3", []byte{7, 8, 9})

	f.Fuzz(func(t *testing.T, key string, value []byte) {
		field := Bytes(key, value)

		enc := &encoderMock{}
		enc.On("Bytes", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzInt8(f *testing.F) {
	f.Add("key1", int8(16))
	f.Add("key2", int8(23))
	f.Add("key3", int8(63))

	f.Fuzz(func(t *testing.T, key string, value int8) {
		field := Int8(key, value)

		enc := &encoderMock{}
		enc.On("Int8", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzInt16(f *testing.F) {
	f.Add("key1", int16(65))
	f.Add("key2", int16(3))
	f.Add("key3", int16(72))

	f.Fuzz(func(t *testing.T, key string, value int16) {
		field := Int16(key, value)

		enc := &encoderMock{}
		enc.On("Int16", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzInt32(f *testing.F) {
	f.Add("key1", int32(15))
	f.Add("key2", int32(23))
	f.Add("key3", int32(2))

	f.Fuzz(func(t *testing.T, key string, value int32) {
		field := Int32(key, value)

		enc := &encoderMock{}
		enc.On("Int32", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzInt64(f *testing.F) {
	f.Add("key1", int64(875))
	f.Add("key2", int64(85))
	f.Add("key3", int64(54))

	f.Fuzz(func(t *testing.T, key string, value int64) {
		field := Int64(key, value)

		enc := &encoderMock{}
		enc.On("Int64", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzUint16(f *testing.F) {
	f.Add("key1", uint16(16))
	f.Add("key2", uint16(23))
	f.Add("key3", uint16(63))

	f.Fuzz(func(t *testing.T, key string, value uint16) {
		field := Uint16(key, value)

		enc := &encoderMock{}
		enc.On("Uint16", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzUint32(f *testing.F) {
	f.Add("key1", uint32(16))
	f.Add("key2", uint32(23))
	f.Add("key3", uint32(63))

	f.Fuzz(func(t *testing.T, key string, value uint32) {
		field := Uint32(key, value)

		enc := &encoderMock{}
		enc.On("Uint32", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzUint64(f *testing.F) {
	f.Add("key1", uint64(16))
	f.Add("key2", uint64(23))
	f.Add("key3", uint64(63))

	f.Fuzz(func(t *testing.T, key string, value uint64) {
		field := Uint64(key, value)

		enc := &encoderMock{}
		enc.On("Uint64", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzFloat32(f *testing.F) {
	f.Add("key1", float32(16.2))
	f.Add("key2", float32(2.3))
	f.Add("key3", float32(1.63))

	f.Fuzz(func(t *testing.T, key string, value float32) {
		field := Float32(key, value)

		enc := &encoderMock{}
		enc.On("Float32", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func FuzzFloat64(f *testing.F) {
	f.Add("key1", float64(16.2))
	f.Add("key2", float64(2.3))
	f.Add("key3", float64(1.63))

	f.Fuzz(func(t *testing.T, key string, value float64) {
		field := Float64(key, value)

		enc := &encoderMock{}
		enc.On("Float64", key, value)

		field.Encode(enc)
		enc.AssertExpectations(t)
	})
}

func TestTime(t *testing.T) {
	tests := []struct {
		key   string
		value time.Time
	}{
		{
			key:   "key1",
			value: time.Now().Add(5 * time.Minute),
		},
		{
			key:   "key2",
			value: time.Now().Add(10 * time.Minute),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.key, func(t *testing.T) {
			field := Time(test.key, test.value)

			enc := &encoderMock{}
			enc.On("Time", test.key, test.value)

			field.Encode(enc)
			enc.AssertExpectations(t)
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		key   string
		value time.Duration
	}{
		{
			key:   "key1",
			value: 5 * time.Minute,
		},
		{
			key:   "key2",
			value: 10 * time.Minute,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.key, func(t *testing.T) {
			field := Duration(test.key, test.value)

			enc := &encoderMock{}
			enc.On("Duration", test.key, test.value)

			field.Encode(enc)
			enc.AssertExpectations(t)
		})
	}
}

func TestError(t *testing.T) {
	tests := []error{
		errors.New("error1"),
		errors.New("error2"),
		errors.New("error3"),
	}

	for _, test := range tests {
		test := test

		t.Run(test.Error(), func(t *testing.T) {
			field := Error(test)

			enc := &encoderMock{}
			enc.On("Error", test)

			field.Encode(enc)
			enc.AssertExpectations(t)
		})
	}
}
