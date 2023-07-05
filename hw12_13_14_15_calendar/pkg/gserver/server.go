package gserver

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	defaultAddr            = ":8000"
	defaultShutdownTimeout = 20 * time.Second
)

type AppGRPCServer struct {
	GRPC            *grpc.Server
	notify          chan error
	shutdownTimeout time.Duration
	addr            string
}

func NewGrpcServer(grpcServer *grpc.Server, opts ...Option) *AppGRPCServer {
	s := &AppGRPCServer{
		GRPC:            grpcServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
		addr:            defaultAddr,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *AppGRPCServer) GetAddr() string {
	return s.addr
}

func (s *AppGRPCServer) Run() {
	go func() {
		defer close(s.notify)
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			s.notify <- err
			return
		}
		s.notify <- s.GRPC.Serve(lis)
	}()
}

func (s *AppGRPCServer) Notify() <-chan error {
	return s.notify
}

func (s *AppGRPCServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	ok := make(chan struct{})
	go func() {
		s.GRPC.GracefulStop()
		close(ok)
	}()

	select {
	case <-ok:
		return nil
	case <-ctx.Done():
		s.GRPC.Stop()
		return ctx.Err()
	}
}
