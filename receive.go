package cli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	simplemq "github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/message"
)

func runReceiveCommand(ctx context.Context, c *CLI) error {
	logger := slog.With("queue_name", c.Message.QueueName)

	cmd := c.Message.Receive
	client, err := simplemq.NewMessageClient(c.Message.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create message client: %w", err)
	}
	messageOp := simplemq.NewMessageOp(client, c.Message.QueueName)
	logger.Debug("receiving message", "queue", c.Message.QueueName)

	count := 0
	receive := func() error {
		msgs, err := messageOp.Receive(ctx)
		if err != nil {
			return fmt.Errorf("failed to receive message: %w", err)
		}
		for _, msg := range msgs {
			logger.Debug("received message", "message", msg)
			var b []byte
			if cmd.Raw {
				b, err = json.Marshal(msg)
				if err != nil {
					return fmt.Errorf("failed to marshal raw message: %w", err)
				}
			} else {
				m := convertMessageContent(msg, c.Message.Base64)
				b, err = json.Marshal(m)
				if err != nil {
					return fmt.Errorf("failed to marshal message: %w", err)
				}
			}
			fmt.Println(string(b))

			if cmd.AutoDelete {
				logger.Debug("deleting message", "messageID", msg.ID)
				if err := messageOp.Delete(ctx, string(msg.ID)); err != nil {
					return fmt.Errorf("failed to delete message: %w", err)
				}
			}
			count++
			if count >= cmd.Count {
				break
			}
		}
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := receive(); err != nil {
			return err
		}
		if !cmd.Polling {
			break
		}
		if count >= cmd.Count {
			return nil
		}
		logger.Debug("sleeping before next polling", "interval", cmd.Interval)
		sleepWithContext(ctx, cmd.Interval)
	}
	return nil
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return
	case <-t.C:
		return
	}
}

type Message struct {
	ID                  string    `json:"id"`
	Content             string    `json:"content"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	ExpiresAt           time.Time `json:"expires_at"`
	AcquiredAt          time.Time `json:"acquired_at"`
	VisibilityTimeoutAt time.Time `json:"visibility_timeout_at"`
}

func UnixToTime(msec int64) time.Time {
	sec := msec / 1000
	nsec := (msec % 1000) * 1_000_000
	return time.Unix(sec, nsec)
}

func convertMessageContent(msg message.Message, b64 bool) Message {
	m := Message{
		ID:                  string(msg.ID),
		CreatedAt:           UnixToTime(msg.CreatedAt),
		UpdatedAt:           UnixToTime(msg.UpdatedAt),
		ExpiresAt:           UnixToTime(msg.ExpiresAt),
		AcquiredAt:          UnixToTime(msg.AcquiredAt),
		VisibilityTimeoutAt: UnixToTime(msg.VisibilityTimeoutAt),
	}
	if b64 {
		decoded, err := base64.StdEncoding.DecodeString(string(msg.Content))
		if err != nil {
			slog.Warn("failed to decode base64 message content", "error", err, "content", msg.Content)
			m.Content = string(msg.Content)
		} else {
			m.Content = string(decoded)
		}
	} else {
		m.Content = string(msg.Content)
	}
	return m
}
