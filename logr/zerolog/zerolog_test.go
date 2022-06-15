package zerolog

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/blacklane/go-libs/logr"
	"github.com/rs/zerolog"
)

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := New(zerolog.New(&buf))

	logger.Info("test", logr.String("key1", "value1"), logr.Int32("key2", 10))

	expected := `{"level":"info","key1":"value1","key2":10,"message":"test"}`
	assertTrimEqual(t, expected, &buf)
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	logger := New(zerolog.New(&buf))

	logger.Error("test", logr.Bytes("key1", []byte{65, 67}), logr.Error(errors.New("error1")))

	expected := `{"level":"error","key1":"AC","error":"error1","message":"test"}`
	assertTrimEqual(t, expected, &buf)
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := New(zerolog.New(&buf))

	logger = logger.WithFields(
		logr.String("key1", "value1"),
		logr.Error(errors.New("error1")),
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
