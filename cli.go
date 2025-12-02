package cli

import "time"

type CLI struct {
	QueueName string `help:"Queue name" required:"" env:"SIMPLEMQ_QUEUE_NAME"`

	Message *MessageCommand `cmd:"" help:"Message related commands"`
}

type MessageCommand struct {
	APIKey string `help:"API Key" required:"" env:"SIMPLEMQ_API_KEY"`
	Base64 bool   `name:"base64" help:"Use Base64 encoding for message content" default:"false" env:"SIMPLEMQ_BASE64"`

	Send    *SendCommand    `cmd:"" help:"Send message to queue"`
	Receive *ReceiveCommand `cmd:"" help:"Receive message from queue"`
}

type SendCommand struct {
	Content string `arg:"" help:"Content of the message to send" name:"content" env:"SIMPLEMQ_MESSAGE_CONTENT"`
}

type ReceiveCommand struct {
	Polling    bool          `help:"Enable polling to receive message" default:"false" env:"SIMPLEMQ_POLLING"`
	Count      int           `help:"Number of messages to receive" default:"1" env:"SIMPLEMQ_RECEIVE_COUNT"`
	AutoDelete bool          `help:"Automatically delete messages after receiving" default:"false" env:"SIMPLEMQ_AUTO_DELETE"`
	Interval   time.Duration `help:"Polling interval for receiving message" default:"1s" env:"SIMPLEMQ_POLLING_INTERVAL"`
}
