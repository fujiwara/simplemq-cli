package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	simplemq "github.com/sacloud/simplemq-api-go"
)

type ListQueueCommand struct{}

func runListQueueCommand(ctx context.Context, _ *CLI) error {
	// cmd := c.Queue.List
	logger := slog.Default()
	client, err := simplemq.NewQueueClient()
	if err != nil {
		return fmt.Errorf("failed to create queue client: %w", err)
	}
	queueOp := simplemq.NewQueueOp(client)
	logger.Debug("listing queues")
	queues, err := queueOp.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to get message count: %w", err)
	}
	for _, q := range queues {
		b, _ := json.Marshal(q)
		fmt.Println(string(b))
	}
	return nil
}
