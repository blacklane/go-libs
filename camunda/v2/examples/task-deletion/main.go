package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/camunda/v2"
	"github.com/blacklane/go-libs/logger"
)

const (
	url        = "http://localhost:8080"
	processKey = "example-process"
	businessKey = "US_ee14f205-dccc-4d6f-a3ae-5d20a0015d3b"
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

	if err := client.DeleteTaskByBusinessKey(context.Background(), businessKey); err != nil {
		log.Fatal().Err(err).Msg("Failed to Delete Task")
	}
	log.Info().Msg("successful")
}
