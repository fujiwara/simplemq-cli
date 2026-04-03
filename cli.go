package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Version kong.VersionFlag `short:"v" help:"Show version and exit."`
	Debug   bool             `help:"Enable debug mode." env:"SIMPLEMQ_DEBUG" default:"false"`
	Timeout time.Duration    `help:"Timeout for API requests (e.g. 30s, 1m)." env:"SIMPLEMQ_TIMEOUT"`

	Message *MessageCommand `cmd:"" help:"Message related commands"`
	Queue   *QueueCommand   `cmd:"" help:"Queue related commands"`

	w io.Writer `kong:"-"`
	r io.Reader `kong:"-"`
}

func (c *CLI) stdinReader() io.Reader {
	if c.r != nil {
		return c.r
	}
	return os.Stdin
}

type MessageCommand struct {
	QueueName string `name:"queue" help:"Queue name" short:"q" required:"" env:"SIMPLEMQ_QUEUE_NAME"`

	APIKey        string `help:"API Key" required:"" env:"SIMPLEMQ_API_KEY"`
	MessageAPIURL string `help:"Message API URL" env:"SIMPLEMQ_MESSAGE_API_URL"`
	Raw           bool   `help:"Handle raw message without Base64 encoding/decoding" default:"false" env:"SIMPLEMQ_RAW"`

	Send    *SendMessageCommand    `cmd:"" help:"Send message to queue"`
	Receive *ReceiveMessageCommand `cmd:"" help:"Receive message from queue"`
	Delete  *DeleteMessageCommand  `cmd:"" help:"Delete message from queue"`
}

type QueueCommand struct {
	Create       *CreateQueueCommand  `cmd:"" help:"Create a new queue"`
	List         *ListQueueCommand    `cmd:"" help:"List queues"`
	Get          *GetQueueCommand     `cmd:"" help:"Get queue details"`
	Modify       *ModifyQueueCommand  `cmd:"" help:"Modify queue settings"`
	Delete       *DeleteQueueCommand  `cmd:"" help:"Delete a queue"`
	MessageCount *MessageCountCommand `cmd:"" help:"Get message count in a queue"`
	RotateAPIKey *RotateAPIKeyCommand `cmd:"" help:"Rotate API key for a queue"`
	Purge        *PurgeQueueCommand   `cmd:"" help:"Purge all messages in a queue"`
}

type QueueCommandBase struct {
	QueueName string `name:"queue" help:"Queue name" short:"q" env:"SIMPLEMQ_QUEUE_NAME"`
	QueueID   string `name:"queue-id" help:"Queue ID (skip name resolution if specified)" env:"SIMPLEMQ_QUEUE_ID"`
}

func (b *QueueCommandBase) Validate() error {
	if b.QueueName == "" && b.QueueID == "" {
		return fmt.Errorf("either --queue or --queue-id must be specified")
	}
	return nil
}

type ConfirmationCommandBase struct {
	Force bool `help:"Force operation without confirmation prompt" short:"f" default:"false" env:"SIMPLEMQ_FORCE"`
}
