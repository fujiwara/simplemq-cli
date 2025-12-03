package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Songmu/prompter"
	simplemq "github.com/sacloud/simplemq-api-go"
)

type PurgeQueueCommand struct {
	QueueCommandBase
	ConfirmationCommandBase
}

func runPurgeQueueCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.Purge
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

	if !cmd.Force {
		if !prompter.YesNo(fmt.Sprintf("Are you sure you want to purge all messages in the queue '%s'?", cmd.QueueName), false) {
			logger.Info("purge operation cancelled by user")
			return nil
		}
	}

	queueOp := simplemq.NewQueueOp(client)
	logger.Debug("purging queue")
	if err := queueOp.ClearMessages(ctx, simplemq.GetQueueID(queue)); err != nil {
		return fmt.Errorf("failed to purge queue: %w", err)
	}
	logger.Debug("queue purged successfully")
	return nil
}
