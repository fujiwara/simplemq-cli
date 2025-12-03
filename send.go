package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"os"

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

	var rawContent []byte
	if cmd.Content == "-" { // read from stdin
		rawContent, err = readInput(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else {
		rawContent = []byte(cmd.Content)
	}
	var content string
	if c.Message.Raw {
		content = string(rawContent)
	} else {
		// automatic base64 encode
		content = base64.StdEncoding.EncodeToString(rawContent)
	}
	logger.Debug("sending message", "content", content)
	res, err := messageOp.Send(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	logger.Debug("message sent successfully", "messageID", res.ID)
	return nil
}

func readInput(r io.Reader) ([]byte, error) {
	b := new(bytes.Buffer)
	_, err := io.Copy(b, r)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
