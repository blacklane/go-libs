package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/hashicorp/go-multierror"
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

func (g *Graceful) beforeStart() error {
	ctx, cancel := context.WithCancel(g.ctx)
	defer cancel()

	mg := new(multierror.Group)

	for _, fn := range g.opts.beforeStart {
		fn := fn
		mg.Go(func() error {
			if err := fn(ctx); err != nil {
				cancel()
				return err
			}
			return nil
		})
	}
	if err := mg.Wait(); err != nil {
		return err
	}
	return nil
}

func (g *Graceful) afterStop() error {
	mg := new(multierror.Group)
	for _, fn := range g.opts.afterStop {
		fn := fn
		mg.Go(func() error {
			return fn(g.opts.ctx)
		})
	}
	if err := mg.Wait(); err != nil {
		return err
	}
	return nil
}

func (g *Graceful) Run() (gerr error) {
	if g.running.Swap(true) {
		return ErrServersRunning
	}
	defer g.running.Store(false)

	defer func() {
		if err := g.afterStop(); err != nil {
			gerr = multierror.Append(gerr, err)
		}
	}()

	if err := g.beforeStart(); err != nil {
		return err
	}

	tg := new(multierror.Group)
	taskCtx, cancelTasks := context.WithCancel(g.ctx)
	defer cancelTasks()

	for _, task := range g.opts.tasks {
		task := task
		tg.Go(func() error {
			<-taskCtx.Done()
			stopCtx, cancel := context.WithTimeout(g.opts.ctx, g.opts.stopTimeout)
			defer cancel()
			return task.Stop(stopCtx)
		})
		tg.Go(func() error {
			if err := task.Start(taskCtx); err != nil {
				cancelTasks()
				return err
			}
			return nil
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, g.opts.sigs...)

	tg.Go(func() error {
		select {
		case <-taskCtx.Done():
			return nil
		case <-c:
			return g.Stop()
		}
	})

	if err := tg.Wait(); err != nil {
		return err
	}

	return nil
}

func (g *Graceful) Stop() error {
	g.cancel()
	return nil
}
