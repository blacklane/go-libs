package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

var ErrServersRunning = errors.New("servers already running")

type Graceful struct {
	opts    *options
	ctx     context.Context
	cancel  func()
	running atomic.Bool
}

func New(opts ...Option) *Graceful {
	opt := newOptions(opts)
	ctx, cancel := context.WithCancel(opt.ctx)

	return &Graceful{
		ctx:    ctx,
		cancel: cancel,
		opts:   opt,
	}
}

func (g *Graceful) Run() error {
	if g.running.Swap(true) {
		return ErrServersRunning
	}
	defer g.running.Store(false)

	eg, ctx := errgroup.WithContext(g.ctx)
	for _, srv := range g.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx, cancel := context.WithTimeout(g.opts.ctx, g.opts.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		eg.Go(func() error {
			return srv.Start(ctx)
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, g.opts.sigs...)

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return g.Stop()
		}
	})

	if err := eg.Wait(); !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (g *Graceful) Stop() error {
	g.cancel()
	return nil
}
