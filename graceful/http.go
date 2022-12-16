package graceful

import (
	"context"
	"errors"
	"net/http"
)

type httpServerTask struct {
	server *http.Server
}

func NewHTTPServerTask(server *http.Server) Task {
	return &httpServerTask{
		server: server,
	}
}

func (s *httpServerTask) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *httpServerTask) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
