package cli

import (
	"context"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
)

type MessageCountCommand struct {
	QueueCommandBase
}

func runMessageCountCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.MessageCount
	logger := slog.With("queue_name", cmd.QueueName)
	client, err := simplemq.NewQueueClient(&queueClient)
	if err != nil {
		return fmt.Errorf("failed to create queue client: %w", err)
	}
	queue, err := resolveQueueByBase(ctx, client, &cmd.QueueCommandBase)
	if err != nil {
		return fmt.Errorf("failed to get queue details: %w", err)
	}
	queueOp := simplemq.NewQueueOp(client)
	logger.Debug("getting message count")
	count, err := queueOp.CountMessages(ctx, simplemq.GetQueueID(queue))
	if err != nil {
		return fmt.Errorf("failed to get message count: %w", err)
	}
	logger.Debug("message count retrieved successfully", "count", count)
	fmt.Fprintln(c.w, count)
	return nil
}
