package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/fujiwara/simplemq-cli/localserver"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), signals()...)
	defer stop()
	if err := run(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var addr string
	var apiKey string
	flag.StringVar(&addr, "addr", "127.0.0.1:18080", "listen address")
	flag.StringVar(&apiKey, "api-key", "", "API key for authentication (if empty, any key is accepted)")
	flag.Parse()

	handler := localserver.NewHandler(apiKey)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	fmt.Fprintf(os.Stderr, "simplemq-localserver listening on %s\n", addr)
	fmt.Fprintf(os.Stderr, "\nTo use with simplemq-cli:\n")
	fmt.Fprintf(os.Stderr, "  export SIMPLEMQ_MESSAGE_API_URL=http://%s\n", addr)
	if apiKey != "" {
		fmt.Fprintf(os.Stderr, "  export SIMPLEMQ_API_KEY=%s\n\n", apiKey)
	} else {
		fmt.Fprintf(os.Stderr, "  export SIMPLEMQ_API_KEY=dummy\n\n")
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
