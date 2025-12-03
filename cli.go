package cli

import (
	"time"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Version kong.VersionFlag `short:"v" help:"Show version and exit."`
	Debug   bool             `help:"Enable debug mode." env:"SIMPLEMQ_DEBUG" default:"false"`

	Message *MessageCommand `cmd:"" help:"Message related commands"`
}

type MessageCommand struct {
	QueueName string `help:"Queue name" required:"" env:"SIMPLEMQ_QUEUE_NAME"`
	APIKey    string `help:"API Key" required:"" env:"SIMPLEMQ_API_KEY"`
	Raw       bool   `help:"Handle raw message without Base64 encoding/decoding" default:"false" env:"SIMPLEMQ_RAW"`

	Send    *SendCommand    `cmd:"" help:"Send message to queue"`
	Receive *ReceiveCommand `cmd:"" help:"Receive message from queue"`
	Delete  *DeleteCommand  `cmd:"" help:"Delete message from queue"`
}

type SendCommand struct {
	Content string `arg:"" help:"Content of the message to send. if - read from stdin" name:"content"`
}

type ReceiveCommand struct {
	Polling    bool          `help:"Enable polling to receive message" default:"false" env:"SIMPLEMQ_POLLING"`
	Count      int           `help:"Number of messages to receive" default:"1" env:"SIMPLEMQ_RECEIVE_COUNT"`
	AutoDelete bool          `help:"Automatically delete messages after receiving" default:"false" env:"SIMPLEMQ_AUTO_DELETE"`
	Interval   time.Duration `help:"Polling interval for receiving message" default:"1s" env:"SIMPLEMQ_POLLING_INTERVAL"`
}

type DeleteCommand struct {
	MessageID string `arg:"" help:"ID of the message to delete" name:"message-id"`
}
