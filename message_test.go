package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/fujiwara/simplemq-cli/localserver"
)

func newTestCLI(serverURL, queueName string) *CLI {
	return &CLI{
		w: &bytes.Buffer{},
		Message: &MessageCommand{
			QueueCommandBase: QueueCommandBase{
				QueueName: queueName,
			},
			APIKey:        "test-api-key",
			MessageAPIURL: serverURL,
		},
	}
}

func outputString(c *CLI) string {
	return c.w.(*bytes.Buffer).String()
}

func resetOutput(c *CLI) {
	c.w.(*bytes.Buffer).Reset()
}

func TestMessageSendAndReceive(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.URL(), "test-queue")
	c.Message.Send = &SendMessageCommand{Content: "hello world"}

	if err := runSendMessageCommand(ctx, c); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	c.Message.Receive = &ReceiveMessageCommand{Count: 1}
	if err := runReceiveMessageCommand(ctx, c); err != nil {
		t.Fatalf("receive failed: %v", err)
	}

	output := outputString(c)
	var msg Message
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &msg); err != nil {
		t.Fatalf("failed to unmarshal output: %v\noutput: %s", err, output)
	}
	if msg.Content != "hello world" {
		t.Errorf("expected content 'hello world', got %q", msg.Content)
	}
	if msg.ID == "" {
		t.Error("expected non-empty message ID")
	}
}

func TestMessageReceiveEmpty(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.URL(), "empty-queue")
	c.Message.Receive = &ReceiveMessageCommand{Count: 1}

	if err := runReceiveMessageCommand(ctx, c); err != nil {
		t.Fatalf("receive failed: %v", err)
	}

	if strings.TrimSpace(outputString(c)) != "" {
		t.Errorf("expected empty output, got %q", outputString(c))
	}
}

func TestMessageSendReceiveDelete(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.URL(), "test-queue")
	c.Message.Send = &SendMessageCommand{Content: "delete me"}

	if err := runSendMessageCommand(ctx, c); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	// Receive to get the message ID
	c.Message.Receive = &ReceiveMessageCommand{Count: 1}
	if err := runReceiveMessageCommand(ctx, c); err != nil {
		t.Fatalf("receive failed: %v", err)
	}

	output := outputString(c)
	var msg Message
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &msg); err != nil {
		t.Fatalf("failed to unmarshal output: %v\noutput: %s", err, output)
	}

	// Delete the message
	c.Message.Delete = &DeleteMessageCommand{MessageID: msg.ID}
	if err := runDeleteMessageCommand(ctx, c); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Delete again should fail (message not found)
	if err := runDeleteMessageCommand(ctx, c); err == nil {
		t.Fatal("expected error on second delete, got nil")
	}
}

func TestMessageReceiveAutoDelete(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.URL(), "test-queue")
	c.Message.Send = &SendMessageCommand{Content: "auto delete me"}

	if err := runSendMessageCommand(ctx, c); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	// Receive with auto-delete
	c.Message.Receive = &ReceiveMessageCommand{Count: 1, AutoDelete: true}
	if err := runReceiveMessageCommand(ctx, c); err != nil {
		t.Fatalf("receive with auto-delete failed: %v", err)
	}

	output := outputString(c)
	var msg Message
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &msg); err != nil {
		t.Fatalf("failed to unmarshal output: %v\noutput: %s", err, output)
	}
	if msg.Content != "auto delete me" {
		t.Errorf("expected content 'auto delete me', got %q", msg.Content)
	}

	// Second receive should return empty (message was auto-deleted)
	resetOutput(c)
	c.Message.Receive = &ReceiveMessageCommand{Count: 1}
	if err := runReceiveMessageCommand(ctx, c); err != nil {
		t.Fatalf("second receive failed: %v", err)
	}

	if strings.TrimSpace(outputString(c)) != "" {
		t.Errorf("expected empty output after auto-delete, got %q", outputString(c))
	}
}
