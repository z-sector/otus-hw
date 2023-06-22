package gserver

import (
	"net"
	"strconv"
	"time"
)

type Option func(*AppGRPCServer)

func Addr(host string, port int) Option {
	return func(s *AppGRPCServer) {
		s.addr = net.JoinHostPort(host, strconv.Itoa(port))
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *AppGRPCServer) {
		s.shutdownTimeout = timeout
	}
}
