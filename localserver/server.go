package localserver

import (
	"net/http"
	"net/http/httptest"
	"time"
)

// Config holds configuration for the local SimpleMQ server.
type Config struct {
	APIKey            string        `help:"API key for authentication (if empty, any key is accepted)" env:"SIMPLEMQ_API_KEY"`
	Addr              string        `help:"Listen address" default:"127.0.0.1:18080" env:"SIMPLEMQ_LOCALSERVER_ADDR"`
	VisibilityTimeout time.Duration `help:"Visibility timeout" default:"30s" env:"SIMPLEMQ_VISIBILITY_TIMEOUT"`
	MessageExpire     time.Duration `help:"Message expire time" default:"96h" env:"SIMPLEMQ_MESSAGE_EXPIRE"`
	Database          string        `help:"SQLite database path for persistent storage (requires sqlite build tag)" env:"SIMPLEMQ_DATABASE"`
	Debug             bool          `help:"Enable debug mode" env:"SIMPLEMQ_DEBUG" default:"false"`
}

// Server is a local SimpleMQ-compatible test server.
type Server struct {
	httpServer *httptest.Server
	mux        *http.ServeMux
	store      Store
	apiKey     string
}

// NewHandler creates a Server as an http.Handler without starting a listener.
// If cfg.APIKey is non-empty, the server validates that incoming requests use this key.
func NewHandler(cfg Config) (*Server, error) {
	store, err := NewStore(cfg.VisibilityTimeout, cfg.MessageExpire, cfg.Database)
	if err != nil {
		return nil, err
	}
	s := &Server{
		store:  store,
		apiKey: cfg.APIKey,
	}
	s.mux = s.buildMux()
	return s, nil
}

// NewServer creates and starts a new local SimpleMQ test server using httptest.
// If cfg.APIKey is non-empty, the server validates that incoming requests use this key.
func NewServer(cfg Config) *Server {
	s, err := NewHandler(cfg)
	if err != nil {
		panic(err)
	}
	s.httpServer = httptest.NewServer(s)
	return s
}

// URL returns the base URL of the test server.
func (s *Server) URL() string {
	return s.httpServer.URL
}

// Close shuts down the test server and closes the store.
func (s *Server) Close() {
	s.httpServer.Close()
	s.store.Close()
}
