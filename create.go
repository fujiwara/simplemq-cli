package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
)

type CreateQueueCommand struct {
	QueueCommandBase
	Description string `help:"Description of the queue"`
}

func runCreateQueueCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.Create
	logger := slog.Default()
	client, err := simplemq.NewQueueClient()
	if err != nil {
		return fmt.Errorf("failed to create queue client: %w", err)
	}
	queueOp := simplemq.NewQueueOp(client)
	logger.Debug("creating queue", "queue_name", cmd.QueueName, "description", cmd.Description)
	q, err := queueOp.Create(ctx, queue.CreateQueueRequest{
		CommonServiceItem: queue.CreateQueueRequestCommonServiceItem{
			Name:        queue.QueueName(cmd.QueueName),
			Description: queue.NewOptString(cmd.Description),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}
	logger.Debug("queue created successfully", "queue", q)
	b, _ := json.Marshal(q)
	fmt.Println(string(b))
	return nil
}
