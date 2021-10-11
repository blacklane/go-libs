package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/x/camunda"
	"github.com/google/uuid"
)

const (
	url        = "http://localhost:8080"
	processKey = "example-process"
	workerID   = "worker-id"
	topic      = "test-topic"
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
	variables := map[string]camunda.CamundaVariable{}
	err := client.StartProcess(context.Background(), businessKey, variables)
	if err != nil {
		log.Err(err).Msg("Failed to start process")
	}

	subscription := client.Subscribe(topic, workerID, func(completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
		log.Info().Msgf("Handling Task [%s] on topic [%s]", t.ID, t.TopicName)

		err := completeFunc(context.Background(), t.ID)
		if err != nil {
			log.Err(err).Msgf("Failed to complete task [%s]", t.ID)
		}
	})

	<-signalChan
	log.Printf("Shutting down")
	subscription.Stop()
}
