package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fujiwara/simplemq-cli/localserver"
)

func newTestCLI(serverURL, queueName string) *CLI {
	return &CLI{
		w: &bytes.Buffer{},
		Message: &MessageCommand{
			QueueName:     queueName,
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
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.TestURL(), "test-queue")
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
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.TestURL(), "empty-queue")
	c.Message.Receive = &ReceiveMessageCommand{Count: 1}

	if err := runReceiveMessageCommand(ctx, c); err != nil {
		t.Fatalf("receive failed: %v", err)
	}

	if strings.TrimSpace(outputString(c)) != "" {
		t.Errorf("expected empty output, got %q", outputString(c))
	}
}

func TestMessageSendReceiveDelete(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.TestURL(), "test-queue")
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

func TestMessageSendStdin(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.TestURL(), "test-stdin")
	c.r = strings.NewReader("hello from stdin")
	c.Message.Send = &SendMessageCommand{Stdin: true}

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
	if msg.Content != "hello from stdin" {
		t.Errorf("expected content 'hello from stdin', got %q", msg.Content)
	}
}

func TestMessageSendEachLine(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.TestURL(), "test-each-line")
	c.r = strings.NewReader("msg1\nmsg2\nmsg3\n")
	c.Message.Send = &SendMessageCommand{Stdin: true, EachLine: true}

	if err := runSendMessageCommand(ctx, c); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	// Receive all 3 messages
	var messages []Message
	for i := 0; i < 3; i++ {
		resetOutput(c)
		c.Message.Receive = &ReceiveMessageCommand{Count: 1, AutoDelete: true}
		if err := runReceiveMessageCommand(ctx, c); err != nil {
			t.Fatalf("receive %d failed: %v", i, err)
		}
		output := strings.TrimSpace(outputString(c))
		if output == "" {
			t.Fatalf("expected message %d, got empty output", i)
		}
		var msg Message
		if err := json.Unmarshal([]byte(output), &msg); err != nil {
			t.Fatalf("failed to unmarshal message %d: %v", i, err)
		}
		messages = append(messages, msg)
	}
	expected := []string{"msg1", "msg2", "msg3"}
	for i, msg := range messages {
		if msg.Content != expected[i] {
			t.Errorf("message[%d] content = %q, want %q", i, msg.Content, expected[i])
		}
	}
}

func TestMessageSendEachJSON(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.TestURL(), "test-each-json")
	c.r = strings.NewReader(`{"a":1} {"b":2}`)
	c.Message.Send = &SendMessageCommand{Stdin: true, EachJSON: true}

	if err := runSendMessageCommand(ctx, c); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	// Receive 2 messages
	var messages []Message
	for i := 0; i < 2; i++ {
		resetOutput(c)
		c.Message.Receive = &ReceiveMessageCommand{Count: 1, AutoDelete: true}
		if err := runReceiveMessageCommand(ctx, c); err != nil {
			t.Fatalf("receive %d failed: %v", i, err)
		}
		output := strings.TrimSpace(outputString(c))
		if output == "" {
			t.Fatalf("expected message %d, got empty output", i)
		}
		var msg Message
		if err := json.Unmarshal([]byte(output), &msg); err != nil {
			t.Fatalf("failed to unmarshal message %d: %v", i, err)
		}
		messages = append(messages, msg)
	}
	if messages[0].Content != `{"a":1}` {
		t.Errorf("message[0] content = %q, want %q", messages[0].Content, `{"a":1}`)
	}
	if messages[1].Content != `{"b":2}` {
		t.Errorf("message[1] content = %q, want %q", messages[1].Content, `{"b":2}`)
	}
}

func TestMessageSendFile(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	// Create a temporary file
	dir := t.TempDir()
	filePath := filepath.Join(dir, "message.txt")
	if err := os.WriteFile(filePath, []byte("hello from file"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	c := newTestCLI(srv.TestURL(), "test-file")
	c.Message.Send = &SendMessageCommand{File: filePath}

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
	if msg.Content != "hello from file" {
		t.Errorf("expected content 'hello from file', got %q", msg.Content)
	}
}

func TestMessageSendFileEachLine(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "lines.txt")
	if err := os.WriteFile(filePath, []byte("line1\nline2\nline3\n"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	c := newTestCLI(srv.TestURL(), "test-file-each-line")
	c.Message.Send = &SendMessageCommand{File: filePath, EachLine: true}

	if err := runSendMessageCommand(ctx, c); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	var messages []Message
	for i := 0; i < 3; i++ {
		resetOutput(c)
		c.Message.Receive = &ReceiveMessageCommand{Count: 1, AutoDelete: true}
		if err := runReceiveMessageCommand(ctx, c); err != nil {
			t.Fatalf("receive %d failed: %v", i, err)
		}
		output := strings.TrimSpace(outputString(c))
		if output == "" {
			t.Fatalf("expected message %d, got empty output", i)
		}
		var msg Message
		if err := json.Unmarshal([]byte(output), &msg); err != nil {
			t.Fatalf("failed to unmarshal message %d: %v", i, err)
		}
		messages = append(messages, msg)
	}
	expected := []string{"line1", "line2", "line3"}
	for i, msg := range messages {
		if msg.Content != expected[i] {
			t.Errorf("message[%d] content = %q, want %q", i, msg.Content, expected[i])
		}
	}
}

func TestMessageSendFileEachJSON(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "data.json")
	if err := os.WriteFile(filePath, []byte(`{"a":1} {"b":2}`), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	c := newTestCLI(srv.TestURL(), "test-file-each-json")
	c.Message.Send = &SendMessageCommand{File: filePath, EachJSON: true}

	if err := runSendMessageCommand(ctx, c); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	var messages []Message
	for i := 0; i < 2; i++ {
		resetOutput(c)
		c.Message.Receive = &ReceiveMessageCommand{Count: 1, AutoDelete: true}
		if err := runReceiveMessageCommand(ctx, c); err != nil {
			t.Fatalf("receive %d failed: %v", i, err)
		}
		output := strings.TrimSpace(outputString(c))
		if output == "" {
			t.Fatalf("expected message %d, got empty output", i)
		}
		var msg Message
		if err := json.Unmarshal([]byte(output), &msg); err != nil {
			t.Fatalf("failed to unmarshal message %d: %v", i, err)
		}
		messages = append(messages, msg)
	}
	if messages[0].Content != `{"a":1}` {
		t.Errorf("message[0] content = %q, want %q", messages[0].Content, `{"a":1}`)
	}
	if messages[1].Content != `{"b":2}` {
		t.Errorf("message[1] content = %q, want %q", messages[1].Content, `{"b":2}`)
	}
}

func TestMessageReceiveAutoDelete(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{})
	defer srv.Close()
	ctx := t.Context()

	c := newTestCLI(srv.TestURL(), "test-queue")
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

func TestTimeout(t *testing.T) {
	srv := localserver.NewTestServer(localserver.Config{Latency: 500 * time.Millisecond})
	defer srv.Close()

	t.Run("timeout exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		c := newTestCLI(srv.TestURL(), "test-timeout")
		c.Message.Send = &SendMessageCommand{Content: "should timeout"}

		err := runSendMessageCommand(ctx, c)
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
		if !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Errorf("expected context deadline exceeded error, got: %v", err)
		}
	})

	t.Run("within timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		c := newTestCLI(srv.TestURL(), "test-timeout")
		c.Message.Send = &SendMessageCommand{Content: "should succeed"}

		if err := runSendMessageCommand(ctx, c); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})
}
