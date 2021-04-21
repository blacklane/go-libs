package main

import (
	"log"

	"github.com/blacklane/go-libs/x/probes"
)

func main() {
	p := probes.New(":4242")

	// Start the probes on a goroutine as p.Start() is a blocking call
	go func() {
		if err := p.Start(); err != nil {
			log.Printf("could not start probes: %v\n", err)
		}
	}()
}
