package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
	queue "github.com/sacloud/simplemq-api-go/apis/v1/queue"
)

type ModifyQueueCommand struct {
	QueueCommandBase
	VisibilityTimeoutSeconds *int64 `help:"Visibility timeout in seconds"`
	ExpireSeconds            *int64 `help:"Message expire time in seconds"`
}

func runModifyQueueCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.Modify
	logger := slog.With("queue_name", cmd.QueueName)
	client, err := simplemq.NewQueueClient()
	if err != nil {
		return fmt.Errorf("failed to create queue client: %w", err)
	}
	logger.Debug("getting queue details")
	q, err := resolveQueue(ctx, client, cmd.QueueName)
	if err != nil {
		return fmt.Errorf("failed to get queue details: %w", err)
	}
	logger.Debug("queue details retrieved successfully", "queue", q)

	if cmd.VisibilityTimeoutSeconds == nil && cmd.ExpireSeconds == nil {
		return fmt.Errorf("no modification parameters provided")
	}

	queueOp := simplemq.NewQueueOp(client)
	queueID := simplemq.GetQueueID(q)

	settings := q.Settings
	if cmd.VisibilityTimeoutSeconds != nil {
		settings.VisibilityTimeoutSeconds = queue.VisibilityTimeoutSeconds(*cmd.VisibilityTimeoutSeconds)
	}
	if cmd.ExpireSeconds != nil {
		settings.ExpireSeconds = queue.ExpireSeconds(*cmd.ExpireSeconds)
	}
	logger.Debug("modifying queue", "settings", settings)

	res, err := queueOp.Config(ctx, queueID, queue.ConfigQueueRequest{
		CommonServiceItem: queue.ConfigQueueRequestCommonServiceItem{
			Settings: settings,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to modify queue: %w", err)
	}
	logger.Debug("queue modified successfully", "queue", res)
	b, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal queue: %w", err)
	}
	fmt.Println(string(b))

	return nil
}
