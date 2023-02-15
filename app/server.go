package app

import (
	"fmt"
	"log"
	"net/http"
)

// ServerOptionFunc can be used to customize an http Server
type ServerOptionFunc func(*Server) error

// Server wraps an http.Server
type Server struct {
	name   string
	server *http.Server
}

// WithAddr sets the address of the http Server to listen on
func WithAddr(addr string) ServerOptionFunc {
	return func(s *Server) error {
		s.server.Addr = addr
		return nil
	}
}

func withHandler(h http.Handler) ServerOptionFunc {
	return func(s *Server) error {
		s.server.Handler = h
		return nil
	}
}

func newServer(prefix string, options ...ServerOptionFunc) (*Server, error) {
	srv := &Server{
		name:   prefix,
		server: &http.Server{},
	}

	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(srv); err != nil {
			return nil, fmt.Errorf("unable to create a new http server instance: %w", err)
		}
	}
	return srv, nil
}

// ListenAndServe starts the http server
func (s *Server) ListenAndServe() {
	log.Printf("listening at %s...\n", s.server.Addr)
	s.server.ListenAndServe()
}
