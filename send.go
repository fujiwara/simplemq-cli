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

	client, err := simplemq.NewMessageClient(c.Message.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create message client: %w", err)
	}
	messageOp := simplemq.NewMessageOp(client, c.QueueName)

	var content string
	if c.Message.Base64 {
		content = base64.StdEncoding.EncodeToString([]byte(cmd.Content))
	} else {
		content = cmd.Content
	}
	slog.Info("sending message", "queue", c.QueueName, "content", content)
	_, err = messageOp.Send(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	slog.Info("message sent successfully")
	return nil
}
