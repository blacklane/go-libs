package zerolog

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/blacklane/go-libs/logr/field"
	"github.com/rs/zerolog"
)

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := New(zerolog.New(&buf))

	logger.Info("test", field.String("key1", "value1"), field.Int32("key2", 10))

	expected := `{"level":"info","key1":"value1","key2":10,"message":"test"}`
	assertTrimEqual(t, expected, &buf)
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	logger := New(zerolog.New(&buf))

	logger.Error(errors.New("error1"), "test", field.Bytes("key1", []byte{65, 67}))

	expected := `{"level":"error","error":"error1","key1":"AC","message":"test"}`
	assertTrimEqual(t, expected, &buf)
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := New(zerolog.New(&buf))

	logger = logger.WithFields(
		field.String("key1", "value1"),
		field.Error(errors.New("error1")),
	)
	logger.Info("test")

	expected := `{"level":"info","key1":"value1","error":"error1","message":"test"}`
	assertTrimEqual(t, expected, &buf)
}

func assertTrimEqual(t *testing.T, expected string, actual fmt.Stringer) {
	if strings.Trim(actual.String(), "\n") != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
