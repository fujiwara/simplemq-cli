package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
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
	var cfg localserver.Config
	kong.Parse(&cfg)

	handler := localserver.NewHandler(cfg)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	fmt.Fprintf(os.Stderr, "simplemq-localserver listening on %s\n", cfg.Addr)
	fmt.Fprintf(os.Stderr, "\nTo use with simplemq-cli:\n")
	fmt.Fprintf(os.Stderr, "  export SIMPLEMQ_MESSAGE_API_URL=http://%s\n", cfg.Addr)
	if cfg.APIKey != "" {
		fmt.Fprintf(os.Stderr, "  export SIMPLEMQ_API_KEY=%s\n\n", cfg.APIKey)
	} else {
		fmt.Fprintf(os.Stderr, "  export SIMPLEMQ_API_KEY=dummy\n\n")
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
