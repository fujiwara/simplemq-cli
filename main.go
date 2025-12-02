package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

func Run(ctx context.Context) error {
	var c CLI
	k, err := kong.New(&c)
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}
	kx, err := k.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}
	switch kx.Command() {
	case "message send <content>":
		return runSendCommand(ctx, &c)
	case "message receive":
		return runReceiveCommand(ctx, &c)
	default:
		return fmt.Errorf("unknown command: %s", kx.Command())
	}
}
