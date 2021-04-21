package probes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const (
	liveResponse  = `{"status":"Kubernetes I'm live, no need to restart me"}`
	readyResponse = `{"status":"Kubernetes I'm ready, you can send requests to me"}`
)

type Probe interface {
	Start() error
	Shutdown(ctx context.Context) error

	// WithContext returns a copy of parent associated with this Probe
	WithContext(parent context.Context) context.Context

	LivenessSucceed()
	LivenessFail()

	ReadinessSucceed()
	ReadinessFail()
}

type probe struct {
	server *http.Server

	live  bool
	ready bool

	liveChan  chan bool
	readyChan chan bool
}

type noop struct{}

// Start is a noop and returns nil
func (n *noop) Start() error { return nil }

// Shutdown is a noop and returns nil
func (n *noop) Shutdown(context.Context) error { return nil }

// WithContext is a noop and returns ctx
func (n *noop) WithContext(ctx context.Context) context.Context { return ctx }

// LivenessSucceed is a noop
func (n *noop) LivenessSucceed() {}

// LivenessFail is a noop
func (n *noop) LivenessFail() {}

// ReadinessSucceed is a noop
func (n *noop) ReadinessSucceed() {}

// ReadinessFail is a noop
func (n *noop) ReadinessFail() {}

// New creates an HTTP server to handle the kubernetes liveness and readiness probes
// on /live and /ready respectively.
// To start te server call Start(), to shut it down, Shutdown()
// addr is the server address check Addr property of http.Server for details.
func New(addr string) Probe {
	p := &probe{
		live:  true,
		ready: true,

		liveChan:  make(chan bool),
		readyChan: make(chan bool),
	}

	mux := chi.NewMux()
	mux.Get("/live", p.liveHandler)
	mux.Get("/ready", p.readyHandler)

	p.server = &http.Server{
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: time.Second,
	}

	return p
}

func NewNoop() Probe {
	return &noop{}
}

// Start starts a HTTP server to serve the liveness and readiness probes.
// It's a blocking call, you might want to run it on a goroutine:
// 		go func() {
// 			if err := p.Start(); err != nil {
// 				log.Printf("could not start probes: %v\n", err)
// 			}
// 		}()
func (p *probe) Start() error {
	go func() {
		for {
			select {
			case live := <-p.liveChan:
				p.live = live
			case ready := <-p.readyChan:
				p.ready = ready
			}
			fmt.Println("chan for loop is looping")
		}
	}()
	return p.server.ListenAndServe()
}

// Shutdown shuts down the associated http.Server
func (p *probe) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *probe) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, p)
}

func (p *probe) LivenessSucceed() { p.liveChan <- true }

func (p *probe) LivenessFail() { p.liveChan <- false }

func (p *probe) ReadinessSucceed() { p.readyChan <- true }

func (p *probe) ReadinessFail() { p.readyChan <- false }

func (p *probe) liveHandler(w http.ResponseWriter, r *http.Request) {
	if !p.live {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error, check the logs"))
		return
	}

	_, _ = w.Write([]byte(liveResponse))
}

func (p *probe) readyHandler(w http.ResponseWriter, r *http.Request) {
	if !p.ready {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error, check the logs"))
		return
	}

	_, _ = w.Write([]byte(readyResponse))
}
