package localserver

import (
	"net/http"
	"net/http/httptest"
)

// Server is a local SimpleMQ-compatible test server.
type Server struct {
	httpServer *httptest.Server
	mux        *http.ServeMux
	store      *Store
}

// NewHandler creates a Server as an http.Handler without starting a listener.
func NewHandler() *Server {
	s := &Server{
		store: NewStore(),
	}
	s.mux = s.buildMux()
	return s
}

// NewServer creates and starts a new local SimpleMQ test server using httptest.
func NewServer() *Server {
	s := NewHandler()
	s.httpServer = httptest.NewServer(s)
	return s
}

// URL returns the base URL of the test server.
func (s *Server) URL() string {
	return s.httpServer.URL
}

// Close shuts down the test server.
func (s *Server) Close() {
	s.httpServer.Close()
}
