package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Songmu/prompter"
	simplemq "github.com/sacloud/simplemq-api-go"
)

type DeleteMessageCommand struct {
	MessageID string `arg:"" help:"ID of the message to delete" name:"message-id"`
}

func runDeleteMessageCommand(ctx context.Context, c *CLI) error {
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

type DeleteQueueCommand struct {
	QueueCommandBase
	ConfirmationCommandBase
}

func runDeleteQueueCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.Delete
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
		if !prompter.YesNo(fmt.Sprintf("Are you sure you want to delete the queue '%s'?", cmd.QueueName), false) {
			logger.Info("delete operation cancelled by user")
			return nil
		}
	}

	queueOp := simplemq.NewQueueOp(client)
	logger.Debug("deleting queue")
	if err := queueOp.Delete(ctx, simplemq.GetQueueID(queue)); err != nil {
		return fmt.Errorf("failed to delete queue: %w", err)
	}
	logger.Debug("queue deleted successfully")
	return nil
}
