package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/blacklane/go-libs/x/camunda"
)

const url = "http://localhost:8080"
const processKey = "example-process"

func main() {
	// catch the signals as soon as possible
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // a.k.a ctrl+C

	// camunda stuff
	client := camunda.NewClient(nil, url, processKey, nil)

	<-signalChan
	log.Printf("Shutting down")
}
