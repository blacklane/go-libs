package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/logger"
	"github.com/google/uuid"

	"github.com/blacklane/go-libs/x/camunda"
)

const url = "http://localhost:8080"
const processKey = "example-process"

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	log := logger.New(
		os.Stdout,
		"camunda-sample-client",
		logger.WithLevel("debug"))

	// camunda stuff
	client := camunda.NewClient(log, url, processKey, nil)

	businessKey := uuid.New().String()
	variables := map[string]camunda.CamundaVariable{}
	err := client.StartProcess(context.Background(), businessKey, variables)
	if err != nil {
		log.Err(err).Msg("Failed to start process")
	}

	<-signalChan
	log.Printf("Shutting down")
}
