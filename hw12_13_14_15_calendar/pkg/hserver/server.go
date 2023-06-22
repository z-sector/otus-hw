package hserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":8000"
	_defaultShutdownTimeout = 5 * time.Second
)

type AppHTTPServer struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewHTTPServer(handler http.Handler, opts ...Option) *AppHTTPServer {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &AppHTTPServer{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *AppHTTPServer) GetAddr() string {
	return s.server.Addr
}

func (s *AppHTTPServer) Run() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *AppHTTPServer) Notify() <-chan error {
	return s.notify
}

func (s *AppHTTPServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
