package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/logger"
	"github.com/google/uuid"

	"github.com/blacklane/go-libs/camunda/v2"
)

const (
	url        = "http://localhost:8080"
	processKey = "example-process"
)

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	log := logger.New(
		os.Stdout,
		"camunda-sample-client",
		logger.WithLevel("debug"))

	credentials := camunda.BasicAuthCredentials{
		User:     "Bernd",
		Password: "password",
	}
	// camunda stuff
	client := camunda.NewClient(url, processKey, http.Client{}, credentials)

	businessKey := uuid.New().String()
	variables := map[string]camunda.Variable{}
	err := client.StartProcess(context.Background(), businessKey, variables)
	if err != nil {
		log.Err(err).Msg("Failed to start process")
	}

	err = client.SendMessage(context.Background(), "set-color", businessKey, map[string]camunda.Variable{
		"color": camunda.NewStringVariable(camunda.VarTypeString, "yellow"),
	})
	if err != nil {
		log.Err(err).Msg("Failed to send message")
	}
}
