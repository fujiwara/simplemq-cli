package cli

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
)

func runSendCommand(ctx context.Context, c *CLI) error {
	cmd := c.Message.Send
	logger := slog.With("queue_name", c.Message.QueueName)

	client, err := simplemq.NewMessageClient(c.Message.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create message client: %w", err)
	}
	messageOp := simplemq.NewMessageOp(client, c.Message.QueueName)

	var content string
	if c.Message.Raw {
		content = cmd.Content
	} else {
		// automatic base64 encode
		content = base64.StdEncoding.EncodeToString([]byte(cmd.Content))
	}
	logger.Debug("sending message", "content", content)
	res, err := messageOp.Send(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	logger.Debug("message sent successfully", "messageID", res.ID)
	return nil
}
