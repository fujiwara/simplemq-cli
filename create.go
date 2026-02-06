package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sacloud/saclient-go"
	simplemq "github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/message"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
)

var queueClient saclient.Client

func newMessageClient(c *CLI) (*message.Client, error) {
	if u := c.Message.MessageAPIURL; u != "" {
		return simplemq.NewMessageClientWithApiUrl(u, c.Message.APIKey, &queueClient)
	}
	return simplemq.NewMessageClient(c.Message.APIKey, &queueClient)
}

type CreateQueueCommand struct {
	QueueCommandBase
	Description string `help:"Description of the queue"`
}

func runCreateQueueCommand(ctx context.Context, c *CLI) error {
	cmd := c.Queue.Create
	logger := slog.Default()
	client, err := simplemq.NewQueueClient(&queueClient)
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
