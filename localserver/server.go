package localserver

import (
	"net/http"
	"net/http/httptest"
	"time"
)

// Config holds configuration for the local SimpleMQ server.
type Config struct {
	APIKey                   string `help:"API key for authentication (if empty, any key is accepted)" env:"SIMPLEMQ_API_KEY"`
	Addr                     string `help:"Listen address" default:"127.0.0.1:18080" env:"SIMPLEMQ_LOCALSERVER_ADDR"`
	VisibilityTimeoutSeconds int    `help:"Visibility timeout in seconds" default:"30" env:"SIMPLEMQ_VISIBILITY_TIMEOUT_SECONDS"`
	MessageExpireSeconds     int    `help:"Message expire time in seconds" default:"345600" env:"SIMPLEMQ_MESSAGE_EXPIRE_SECONDS"`
}

// Server is a local SimpleMQ-compatible test server.
type Server struct {
	httpServer *httptest.Server
	mux        *http.ServeMux
	store      *Store
	apiKey     string
}

// NewHandler creates a Server as an http.Handler without starting a listener.
// If cfg.APIKey is non-empty, the server validates that incoming requests use this key.
func NewHandler(cfg Config) *Server {
	visibilityTimeout := time.Duration(cfg.VisibilityTimeoutSeconds) * time.Second
	messageExpiration := time.Duration(cfg.MessageExpireSeconds) * time.Second
	s := &Server{
		store:  NewStore(visibilityTimeout, messageExpiration),
		apiKey: cfg.APIKey,
	}
	s.mux = s.buildMux()
	return s
}

// NewServer creates and starts a new local SimpleMQ test server using httptest.
// If cfg.APIKey is non-empty, the server validates that incoming requests use this key.
func NewServer(cfg Config) *Server {
	s := NewHandler(cfg)
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
