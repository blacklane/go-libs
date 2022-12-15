package graceful

import "context"

type Task interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type task struct {
	start func(context.Context) error
	stop  func(context.Context) error
}

func NewTask(
	start func(context.Context) error,
	stop func(context.Context) error,
) Task {
	return &task{
		start: start,
		stop:  stop,
	}
}

func (s *task) Start(ctx context.Context) error {
	if s.start != nil {
		return s.start(ctx)
	}
	return nil
}

func (s *task) Stop(ctx context.Context) error {
	if s.stop != nil {
		return s.stop(ctx)
	}
	return nil
}
