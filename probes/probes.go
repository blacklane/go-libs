package probes

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

const (
	liveResponse  = `{"status":"Kubernetes I'm live, no need to restart me"}`
	readyResponse = `{"status":"Kubernetes I'm ready, you can send requests to me"}`
)

//revive:disable-line:exported Obvious interface.
type Probe interface {
	// Start starts the probes.
	Start() error
	// Shutdown shuts the probes down.
	Shutdown(ctx context.Context) error

	// WithContext returns a copy of parent associated with this Probe.
	WithContext(parent context.Context) context.Context

	// LivenessSucceed sets the liveness probe to respond success, 200 - OK.
	LivenessSucceed()
	// LivenessFail sets the liveness probe to return failure, 500 - Internal Server Error.
	LivenessFail()

	// ReadinessSucceed sets the liveness probe to respond success, 200 - OK.
	ReadinessSucceed()
	// ReadinessFail sets the liveness probe to return failure, 500 - Internal Server Error.
	ReadinessFail()
}

type probe struct {
	server *http.Server

	live  bool
	ready bool

	rwmu *sync.RWMutex

	liveChan  chan bool
	readyChan chan bool
}

type noop struct{}

// Start is a noop and returns nil.
func (n *noop) Start() error { return nil }

// Shutdown is a noop and returns nil.
func (n *noop) Shutdown(context.Context) error { return nil }

// WithContext is a noop and returns ctx.
func (n *noop) WithContext(ctx context.Context) context.Context { return ctx }

// LivenessSucceed is a noop.
func (n *noop) LivenessSucceed() {}

// LivenessFail is a noop.
func (n *noop) LivenessFail() {}

// ReadinessSucceed is a noop.
func (n *noop) ReadinessSucceed() {}

// ReadinessFail is a noop.
func (n *noop) ReadinessFail() {}

// New creates an HTTP server to handle the kubernetes liveness and readiness probes
// on /live and /ready respectively.
// To start te server call Start(), to shut it down, Shutdown()
// addr is the server address check Addr property of http.Server for details.
func New(addr string) Probe {
	p := &probe{
		live:  true,
		ready: true,

		rwmu: &sync.RWMutex{},
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

// NewNoop returns a noop Probe.
func NewNoop() Probe {
	return &noop{}
}

// Start starts a HTTP server to serve the liveness and readiness probes.
// It's a blocking call, you might want to run it on a goroutine:
// 		go func() {
// 			if err := p.Start(); err != nil {
// 				log.Printf("probes shut down: %v\n", err)
// 			}
// 		}()
func (p *probe) Start() error {
	return p.server.ListenAndServe()
}

// Shutdown shuts down the associated http.Server
func (p *probe) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *probe) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, p)
}

func (p *probe) LivenessSucceed() {
	p.rwmu.Lock()
	defer p.rwmu.Unlock()
	p.live = true
}

func (p *probe) LivenessFail() {
	p.rwmu.Lock()
	defer p.rwmu.Unlock()
	p.live = false
}

func (p *probe) ReadinessSucceed() {
	p.rwmu.Lock()
	defer p.rwmu.Unlock()
	p.ready = true
}

func (p *probe) ReadinessFail() {
	p.rwmu.Lock()
	defer p.rwmu.Unlock()
	p.ready = false
}

func (p *probe) liveHandler(w http.ResponseWriter, r *http.Request) {
	p.rwmu.RLock()
	live := p.live
	p.rwmu.RUnlock()

	if !live {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error, check the logs"))
		return
	}

	_, _ = w.Write([]byte(liveResponse))
}

func (p *probe) readyHandler(w http.ResponseWriter, r *http.Request) {
	p.rwmu.RLock()
	ready := p.ready
	p.rwmu.RUnlock()

	if !ready {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error, check the logs"))
		return
	}

	_, _ = w.Write([]byte(readyResponse))
}
