package main

import (
	"errors"

	"github.com/blacklane/go-libs/logr"
	"github.com/blacklane/go-libs/logr/field"
	zaplogr "github.com/blacklane/go-libs/logr/zap"
	"go.uber.org/zap"
)

func init() {
	logger, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		panic(err)
	}

	zlogr := zaplogr.New(logger)
	logr.SetLogger(zlogr)
}

func main() {
	logr.Debug("debug")

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
