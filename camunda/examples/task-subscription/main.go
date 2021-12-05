package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/blacklane/go-libs/logger"

	"github.com/blacklane/go-libs/camunda"
)

const (
	url        = "http://localhost:8080"
	processKey = "example-process"
	workerID   = "worker-id"
	topic      = "test-topic"
)

type taskHandler struct{}

var log logger.Logger

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	log := logger.New(
		os.Stdout,
		"camunda-sample-client",
		logger.WithLevel("debug"))

	// camunda stuff
	client := camunda.NewClient(url, processKey, http.Client{}, camunda.BasicAuthCredentials{})

	handler := taskHandler{}
	subscription := client.Subscribe(context.Background(), topic, workerID, &handler,
		camunda.FetchInterval(time.Second*2),
		camunda.MaxTasksFetch(10),
		camunda.LockDuration(2000),
	)

	<-signalChan
	log.Printf("Shutting down")
	subscription.Stop()
}

func (th *taskHandler) Handle(ctx context.Context, completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
	log.Info().Msgf("Handling Task [%s] on topic [%s]", t.ID, t.TopicName)

	err := completeFunc(context.Background(), t.ID)
	if err != nil {
		log.Err(err).Msgf("Failed to complete task [%s]", t.ID)
	}
}
