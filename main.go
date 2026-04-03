package cli

import (
	"context"
	"errors"
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

	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}

	var cmdErr error
	switch kx.Command() {
	case "message send <content>", "message send":
		cmdErr = runSendMessageCommand(ctx, &c)
	case "message receive":
		cmdErr = runReceiveMessageCommand(ctx, &c)
	case "message delete <message-id>":
		cmdErr = runDeleteMessageCommand(ctx, &c)
	case "queue create":
		cmdErr = runCreateQueueCommand(ctx, &c)
	case "queue get":
		cmdErr = runGetQueueCommand(ctx, &c)
	case "queue list":
		cmdErr = runListQueueCommand(ctx, &c)
	case "queue modify":
		cmdErr = runModifyQueueCommand(ctx, &c)
	case "queue delete":
		cmdErr = runDeleteQueueCommand(ctx, &c)
	case "queue purge":
		cmdErr = runPurgeQueueCommand(ctx, &c)
	case "queue rotate-api-key":
		cmdErr = runRotateQueueAPIKeyCommand(ctx, &c)
	case "queue message-count":
		cmdErr = runMessageCountCommand(ctx, &c)
	default:
		return fmt.Errorf("unknown command: %s", kx.Command())
	}
	if cmdErr != nil {
		if ctx.Err() != nil && errors.Is(cmdErr, context.DeadlineExceeded) {
			return fmt.Errorf("request timed out (--timeout %s): %w", c.Timeout, cmdErr)
		}
		if ctx.Err() != nil && errors.Is(cmdErr, context.Canceled) {
			return fmt.Errorf("request canceled: %w", cmdErr)
		}
	}
	return cmdErr
}
