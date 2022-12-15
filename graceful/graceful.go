package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync/atomic"

	"go.uber.org/multierr"
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

func (g *Graceful) goHooks(ctx context.Context, hooks []Hook) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, hook := range hooks {
		hook := hook
		eg.Go(func() error {
			return hook(ctx)
		})
	}
	return eg.Wait()
}

func (g *Graceful) beforeStart() error {
	eg, ctx := errgroup.WithContext(g.ctx)
	for _, fn := range g.opts.beforeStart {
		fn := fn
		eg.Go(func() error {
			return fn(ctx)
		})
	}
	return eg.Wait()
}

func (g *Graceful) afterStop() error {
	eg := new(errgroup.Group) // without context cancel propagation
	for _, fn := range g.opts.afterStop {
		fn := fn
		eg.Go(func() error {
			return fn(g.opts.ctx)
		})
	}
	return eg.Wait()
}

func (g *Graceful) Run() (gerr error) {
	if g.running.Swap(true) {
		return ErrServersRunning
	}
	defer g.running.Store(false)

	defer func() {
		if err := g.afterStop(); err != nil {
			// replace by "errors.Join" when go1.20 is released:
			// https://tip.golang.org/doc/go1.20#errors
			gerr = multierr.Append(gerr, err)
		}
	}()

	if err := g.goHooks(g.ctx, g.opts.beforeStart); err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(g.ctx)

	for _, task := range g.opts.tasks {
		task := task
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx, cancel := context.WithTimeout(g.opts.ctx, g.opts.stopTimeout)
			defer cancel()
			return task.Stop(stopCtx)
		})
		eg.Go(func() error {
			return task.Start(ctx)
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
