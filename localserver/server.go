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
	apiKey     string
}

// NewHandler creates a Server as an http.Handler without starting a listener.
// If apiKey is non-empty, the server validates that incoming requests use this key.
func NewHandler(apiKey string) *Server {
	s := &Server{
		store:  NewStore(),
		apiKey: apiKey,
	}
	s.mux = s.buildMux()
	return s
}

// NewServer creates and starts a new local SimpleMQ test server using httptest.
// If apiKey is non-empty, the server validates that incoming requests use this key.
func NewServer(apiKey string) *Server {
	s := NewHandler(apiKey)
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
