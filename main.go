package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fujiwara/sloghandler"
)

func Run(ctx context.Context) error {
	var c CLI
	k, err := kong.New(&c, kong.Vars{"version": fmt.Sprintf("simplemq-cli %s", Version)})
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}

	kx, err := k.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}
	c.w = os.Stdout
	level := slog.LevelInfo
	if c.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(sloghandler.NewLogHandler(os.Stderr, &sloghandler.HandlerOptions{
		HandlerOptions: slog.HandlerOptions{Level: level},
		Color:          true,
	})))

	switch kx.Command() {
	case "message send <content>":
		return runSendMessageCommand(ctx, &c)
	case "message receive":
		return runReceiveMessageCommand(ctx, &c)
	case "message delete <message-id>":
		return runDeleteMessageCommand(ctx, &c)
	case "queue create":
		return runCreateQueueCommand(ctx, &c)
	case "queue get":
		return runGetQueueCommand(ctx, &c)
	case "queue list":
		return runListQueueCommand(ctx, &c)
	case "queue modify":
		return runModifyQueueCommand(ctx, &c)
	case "queue delete":
		return runDeleteQueueCommand(ctx, &c)
	case "queue purge":
		return runPurgeQueueCommand(ctx, &c)
	case "queue rotate-api-key":
		return runRotateQueueAPIKeyCommand(ctx, &c)
	case "queue message-count":
		return runMessageCountCommand(ctx, &c)
	default:
		return fmt.Errorf("unknown command: %s", kx.Command())
	}
}
