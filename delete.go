package cli

import (
	"context"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
)

func runDeleteCommand(ctx context.Context, c *CLI) error {
	cmd := c.Message.Delete
	logger := slog.With("queue_name", c.Message.QueueName)

	client, err := simplemq.NewMessageClient(c.Message.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create message client: %w", err)
	}
	messageOp := simplemq.NewMessageOp(client, c.Message.QueueName)

	logger.Debug("deleting message", "messageID", cmd.MessageID)
	if err := messageOp.Delete(ctx, cmd.MessageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	logger.Debug("message deleted successfully")
	return nil
}
