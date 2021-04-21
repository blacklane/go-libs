package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/blacklane/go-libs/x/probes"
)

func main() {
	// Addr is set as the http.Server Addr.
	// See net.Dial for details of the address format.
	p := probes.New(":4242")

	// Start the probes on a goroutine as p.Start() is a blocking call
	go func() {
		if err := p.Start(); err != nil {
			log.Printf("could not start probes: %v\n", err)
		}
	}()

	// Passing it around in the Context
	ctx := p.WithContext(context.Background())
	// Retrieving it from a Context
	pp := probes.FromCtx(ctx)

	// Calls to /ready will fail with HTTP 500.
	// Kubernetes stops sending requests to the pod
	pp.ReadinessFail()
	// Calls to /ready will succeed with HTTP 200 - OK.
	// Kubernetes starts sending requests to the pod again
	pp.ReadinessSucceed()

	// Calls to /live will fail with HTTP 500.
	// Kubernetes will restart the pod
	pp.LivenessFail()
	// Calls to /live will succeed with HTTP 200 - OK.
	// Kubernetes will restart the pod
	pp.LivenessSucceed()

	// Gracefully shutting down the probes HTTP server
	if err := p.Shutdown(context.Background()); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("failed to shutdown probes HTTP server: %v\n", err)
	}

	noop := probes.NewNoop()
}
