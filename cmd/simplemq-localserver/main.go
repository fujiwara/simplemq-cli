package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/fujiwara/simplemq-cli/localserver"
	"github.com/fujiwara/sloghandler"
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

	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(sloghandler.NewLogHandler(os.Stderr, &sloghandler.HandlerOptions{
		HandlerOptions: slog.HandlerOptions{Level: level},
	})))

	handler := localserver.NewHandler(cfg)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	slog.Info("simplemq-localserver starting", "addr", cfg.Addr)
	slog.Info("to use with simplemq-cli",
		"SIMPLEMQ_MESSAGE_API_URL", "http://"+cfg.Addr,
		"SIMPLEMQ_API_KEY", apiKeyHint(cfg.APIKey),
	)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func apiKeyHint(key string) string {
	if key == "" {
		return "dummy (any non-empty value accepted)"
	}
	return "(configured, use the key you specified)"
}
