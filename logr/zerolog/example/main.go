package main

import (
	"errors"

	"github.com/blacklane/go-libs/logr"
	"github.com/blacklane/go-libs/logr/field"
	"github.com/blacklane/go-libs/logr/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	logger := zerolog.New(log.Logger)
	logr.SetLogger(logger)
}

func main() {

	logr.Info("message",
		field.String("key1", "value1"),
		field.Int32("key2", 10),
	)

	logger := logr.WithFields(
		field.Application("my-app"),
		field.IP("127.0.0.1"),
	)

	err := errors.New("error1")

	logger.Error(err, "message", field.Event("my-event"))
}
