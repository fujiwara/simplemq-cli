package cli

import (
	"context"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
)

type RotateAPIKeyCommand struct {
	QueueCommandBase
}

func runRotateQueueAPIKeyCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.RotateAPIKey
	logger := slog.With("queue_name", cmd.QueueName)
	client, err := simplemq.NewQueueClient()
	if err != nil {
		return fmt.Errorf("failed to create queue client: %w", err)
	}
	logger.Debug("getting queue details")
	queue, err := resolveQueue(ctx, client, cmd.QueueName)
	if err != nil {
		return fmt.Errorf("failed to get queue details: %w", err)
	}
	logger.Debug("queue details retrieved successfully", "queue", queue)

	queueOp := simplemq.NewQueueOp(client)
	logger.Debug("rotating API key for queue")
	res, err := queueOp.RotateAPIKey(ctx, simplemq.GetQueueID(queue))
	if err != nil {
		return fmt.Errorf("failed to rotate API key: %w", err)
	}
	logger.Debug("API key rotated successfully")
	fmt.Println(res)
	return nil
}
