package hserver

import (
	"net"
	"strconv"
	"time"
)

type Option func(*AppHTTPServer)

func Addr(host string, port int) Option {
	return func(s *AppHTTPServer) {
		s.server.Addr = net.JoinHostPort(host, strconv.Itoa(port))
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *AppHTTPServer) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *AppHTTPServer) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *AppHTTPServer) {
		s.shutdownTimeout = timeout
	}
}
