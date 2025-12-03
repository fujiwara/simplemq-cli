package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
)

type GetQueueCommand struct {
	QueueCommandBase
}

func runGetQueueCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.Get
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
	b, _ := json.Marshal(queue)
	fmt.Println(string(b))
	return nil
}

func resolveQueue(ctx context.Context, client *queue.Client, name string) (*queue.CommonServiceItem, error) {
	queueOp := simplemq.NewQueueOp(client)
	queues, err := queueOp.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list queues: %w", err)
	}
	for _, q := range queues {
		if q.Name == name {
			return &q, nil
		}
	}
	return nil, fmt.Errorf("queue not found: %s", name)
}
