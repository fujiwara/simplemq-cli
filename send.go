package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	simplemq "github.com/sacloud/simplemq-api-go"
)

type SendMessageCommand struct {
	Content  string `arg:"" optional:"" help:"Content of the message to send. if - read from stdin" name:"content"`
	Stdin    bool   `help:"Read message content from stdin" default:"false" xor:"input"`
	File     string `help:"Read message content from a file" name:"file" type:"existingfile" xor:"input"`
	EachLine bool   `help:"Send each line as a separate message (requires --stdin or --file)" default:"false" name:"each-line"`
	EachJSON bool   `help:"Send each JSON value as a separate message (requires --stdin or --file)" default:"false" name:"each-json"`
}

func (cmd *SendMessageCommand) hasInput() bool {
	return cmd.Stdin || cmd.File != ""
}

func (cmd *SendMessageCommand) Validate() error {
	if cmd.EachLine && !cmd.hasInput() {
		return fmt.Errorf("--each-line requires --stdin or --file")
	}
	if cmd.EachJSON && !cmd.hasInput() {
		return fmt.Errorf("--each-json requires --stdin or --file")
	}
	if cmd.EachLine && cmd.EachJSON {
		return fmt.Errorf("--each-line and --each-json are mutually exclusive")
	}
	if !cmd.hasInput() && cmd.Content == "" {
		return fmt.Errorf("<content> argument, --stdin, or --file is required")
	}
	if cmd.hasInput() && cmd.Content != "" {
		return fmt.Errorf("--stdin/--file and <content> argument are mutually exclusive")
	}
	return nil
}

func runSendMessageCommand(ctx context.Context, c *CLI) error {
	cmd := c.Message.Send
	logger := slog.With("queue_name", c.Message.QueueName)

	client, err := newMessageClient(c)
	if err != nil {
		return fmt.Errorf("failed to create message client: %w", err)
	}
	messageOp := simplemq.NewMessageOp(client, c.Message.QueueName)

	sendOne := func(rawContent []byte) error {
		var content string
		if c.Message.Raw {
			content = string(rawContent)
		} else {
			content = base64.StdEncoding.EncodeToString(rawContent)
		}
		logger.Debug("sending message", "content", content)
		res, err := messageOp.Send(ctx, content)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		logger.Debug("message sent successfully", "messageID", res.ID)
		return nil
	}

	var r io.Reader
	if cmd.File != "" {
		f, err := os.Open(cmd.File)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", cmd.File, err)
		}
		defer f.Close()
		r = f
	} else {
		r = c.stdinReader()
	}

	switch {
	case cmd.EachLine:
		return sendEachLine(r, sendOne)
	case cmd.EachJSON:
		return sendEachJSON(r, sendOne)
	case cmd.Stdin || cmd.File != "" || cmd.Content == "-":
		rawContent, err := readInput(r)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		return sendOne(rawContent)
	default:
		return sendOne([]byte(cmd.Content))
	}
}

func sendEachLine(r io.Reader, sendOne func([]byte) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if err := sendOne(line); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func sendEachJSON(r io.Reader, sendOne func([]byte) error) error {
	dec := json.NewDecoder(r)
	for dec.More() {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return fmt.Errorf("failed to decode JSON: %w", err)
		}
		if err := sendOne(raw); err != nil {
			return err
		}
	}
	return nil
}

func readInput(r io.Reader) ([]byte, error) {
	b := new(bytes.Buffer)
	_, err := io.Copy(b, r)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
