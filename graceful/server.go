package graceful

import "context"

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type server struct {
	start func(context.Context) error
	stop  func(context.Context) error
}

func NewServer(
	start func(context.Context) error,
	stop func(context.Context) error,
) Server {
	return &server{
		start: start,
		stop:  stop,
	}
}

func (s *server) Start(ctx context.Context) error {
	if s.start != nil {
		return s.start(ctx)
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	if s.stop != nil {
		return s.stop(ctx)
	}
	return nil
}
