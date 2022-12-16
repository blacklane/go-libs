package graceful

import (
	"context"
	"os"
	"syscall"
	"time"
)

type Option func(*options)

type Hook func(context.Context) error

type options struct {
	ctx         context.Context
	stopTimeout time.Duration
	tasks       []Task
	sigs        []os.Signal

	// hooks
	beforeStart []Hook
	afterStop   []Hook
}

func newOptions(opts []Option) *options {
	opt := &options{
		ctx:         context.Background(),
		sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		stopTimeout: 15 * time.Second,
	}
	for _, optFn := range opts {
		optFn(opt)
	}
	return opt
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithSignal(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

func WithStopTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.stopTimeout = timeout
	}
}

func WithTasks(servers ...Task) Option {
	return func(o *options) {
		o.tasks = servers
	}
}

func WithBeforeStartHooks(beforeStart ...Hook) Option {
	return func(o *options) {
		o.beforeStart = beforeStart
	}
}

func WithAfterStopHooks(afterStop ...Hook) Option {
	return func(o *options) {
		o.afterStop = afterStop
	}
}
