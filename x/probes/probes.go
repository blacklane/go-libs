package probes

import (
	"context"
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

	LivenessSucceed()
	LivenessFail()

	ReadinessSucceed()
	ReadinessFail()
}

type probe struct {
	server *http.Server

	live  bool
	ready bool
}

// New creates an HTTP server to handle the kubernetes liveness and readiness probes
// on /live and /ready respectively.
// To start te server call Start(), to shut it down, Shutdown()
func New(addr string, opts ...func(*probe)) Probe {
	p := &probe{
		live:  true,
		ready: true,
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

// Start starts a HTTP server to serve the liveness and readiness probes.
// It's a blocking call, you might want to run it on a goroutine.
func (p *probe) Start() error {
	return p.server.ListenAndServe()
}

func (p *probe) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *probe) LivenessSucceed() { p.live = true }

func (p *probe) LivenessFail() { p.live = false }

func (p *probe) ReadinessSucceed() { p.ready = true }

func (p *probe) ReadinessFail() { p.ready = false }

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
