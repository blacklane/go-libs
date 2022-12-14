package graceful

import (
	"context"
	"errors"
	"net/http"
)

type httpServer struct {
	server *http.Server
}

func NewHTTPServer(server *http.Server) Server {
	return &httpServer{
		server: server,
	}
}

func (s *httpServer) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *httpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
