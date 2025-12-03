package cli

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	Version kong.VersionFlag `short:"v" help:"Show version and exit."`
	Debug   bool             `help:"Enable debug mode." env:"SIMPLEMQ_DEBUG" default:"false"`

	Message *MessageCommand `cmd:"" help:"Message related commands"`
	Queue   *QueueCommand   `cmd:"" help:"Queue related commands"`
}

type MessageCommand struct {
	QueueCommandBase

	APIKey string `help:"API Key" required:"" env:"SIMPLEMQ_API_KEY"`
	Raw    bool   `help:"Handle raw message without Base64 encoding/decoding" default:"false" env:"SIMPLEMQ_RAW"`

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
	QueueName string `help:"Queue name" short:"q" required:"" env:"SIMPLEMQ_QUEUE_NAME"`
}

type ConfirmationCommandBase struct {
	Force bool `help:"Force operation without confirmation prompt" short:"f" default:"false" env:"SIMPLEMQ_FORCE"`
}
