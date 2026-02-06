package localserver_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/fujiwara/simplemq-cli/localserver"
	"github.com/sacloud/simplemq-api-go/apis/v1/message"
)

type testSecuritySource struct {
	token string
}

func (s *testSecuritySource) ApiKeyAuth(_ context.Context, _ message.OperationName) (message.ApiKeyAuth, error) {
	return message.ApiKeyAuth{Token: s.token}, nil
}

func newTestClient(t *testing.T, serverURL, token string) *message.Client {
	t.Helper()
	client, err := message.NewClient(serverURL, &testSecuritySource{token: token})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func b64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// nonexistentUUID is a valid UUID format that will never match a real message ID.
const nonexistentUUID = "00000000-0000-0000-0000-000000000000"

func TestSendReceiveDelete(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()

	client := newTestClient(t, srv.URL(), "test-api-key")
	ctx := context.Background()
	queueName := "test-queue"
	content := b64("hello")

	// Send a message
	sendRes, err := client.SendMessage(ctx, &message.SendRequest{Content: message.MessageContent(content)}, message.SendMessageParams{QueueName: message.QueueName(queueName)})
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}
	sendOK, ok := sendRes.(*message.SendMessageOK)
	if !ok {
		t.Fatalf("expected SendMessageOK, got %T", sendRes)
	}
	if sendOK.Result != "success" {
		t.Errorf("expected result=success, got %s", sendOK.Result)
	}
	if string(sendOK.Message.Content) != content {
		t.Errorf("expected content=%s, got %s", content, sendOK.Message.Content)
	}
	if sendOK.Message.ID == "" {
		t.Error("expected non-empty message ID")
	}
	msgID := sendOK.Message.ID

	// Receive the message
	recvRes, err := client.ReceiveMessage(ctx, message.ReceiveMessageParams{QueueName: message.QueueName(queueName)})
	if err != nil {
		t.Fatalf("receive failed: %v", err)
	}
	recvOK, ok := recvRes.(*message.ReceiveMessageOK)
	if !ok {
		t.Fatalf("expected ReceiveMessageOK, got %T", recvRes)
	}
	if recvOK.Result != "success" {
		t.Errorf("expected result=success, got %s", recvOK.Result)
	}
	if len(recvOK.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(recvOK.Messages))
	}
	recvMsg := recvOK.Messages[0]
	if recvMsg.ID != msgID {
		t.Errorf("expected message ID=%s, got %s", msgID, recvMsg.ID)
	}
	if string(recvMsg.Content) != content {
		t.Errorf("expected content=%s, got %s", content, recvMsg.Content)
	}
	if recvMsg.AcquiredAt == 0 {
		t.Error("expected non-zero acquired_at")
	}
	if recvMsg.VisibilityTimeoutAt == 0 {
		t.Error("expected non-zero visibility_timeout_at")
	}

	// Delete the message
	delRes, err := client.DeleteMessage(ctx, message.DeleteMessageParams{
		QueueName: message.QueueName(queueName),
		MessageId: msgID,
	})
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	delOK, ok := delRes.(*message.DeleteMessageOK)
	if !ok {
		t.Fatalf("expected DeleteMessageOK, got %T", delRes)
	}
	if delOK.Result != "success" {
		t.Errorf("expected result=success, got %s", delOK.Result)
	}
}

func TestEmptyReceive(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()

	client := newTestClient(t, srv.URL(), "test-api-key")
	ctx := context.Background()

	recvRes, err := client.ReceiveMessage(ctx, message.ReceiveMessageParams{QueueName: "empty-queue"})
	if err != nil {
		t.Fatalf("receive failed: %v", err)
	}
	recvOK, ok := recvRes.(*message.ReceiveMessageOK)
	if !ok {
		t.Fatalf("expected ReceiveMessageOK, got %T", recvRes)
	}
	if len(recvOK.Messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(recvOK.Messages))
	}
}

func TestDeleteNotFound(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()

	client := newTestClient(t, srv.URL(), "test-api-key")
	ctx := context.Background()

	delRes, err := client.DeleteMessage(ctx, message.DeleteMessageParams{
		QueueName: "test-queue",
		MessageId: nonexistentUUID,
	})
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	if _, ok := delRes.(*message.DeleteMessageNotFound); !ok {
		t.Fatalf("expected DeleteMessageNotFound, got %T", delRes)
	}
}

func TestUnauthorized(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()

	// Use empty token to trigger 401
	client := newTestClient(t, srv.URL(), "")
	ctx := context.Background()

	sendRes, err := client.SendMessage(ctx, &message.SendRequest{Content: message.MessageContent(b64("hello"))}, message.SendMessageParams{QueueName: "test-queue"})
	if err != nil {
		t.Fatalf("send request failed: %v", err)
	}
	if _, ok := sendRes.(*message.SendMessageUnauthorized); !ok {
		t.Fatalf("expected SendMessageUnauthorized, got %T", sendRes)
	}
}

func TestExtendTimeout(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()

	client := newTestClient(t, srv.URL(), "test-api-key")
	ctx := context.Background()
	queueName := "test-queue"

	// Send and receive a message
	sendRes, err := client.SendMessage(ctx, &message.SendRequest{Content: message.MessageContent(b64("timeout test"))}, message.SendMessageParams{QueueName: message.QueueName(queueName)})
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}
	sendOK := sendRes.(*message.SendMessageOK)
	msgID := sendOK.Message.ID

	_, err = client.ReceiveMessage(ctx, message.ReceiveMessageParams{QueueName: message.QueueName(queueName)})
	if err != nil {
		t.Fatalf("receive failed: %v", err)
	}

	// Extend timeout
	extRes, err := client.ExtendMessageTimeout(ctx, message.ExtendMessageTimeoutParams{
		QueueName: message.QueueName(queueName),
		MessageId: msgID,
	})
	if err != nil {
		t.Fatalf("extend timeout failed: %v", err)
	}
	extOK, ok := extRes.(*message.ExtendMessageTimeoutOK)
	if !ok {
		t.Fatalf("expected ExtendMessageTimeoutOK, got %T", extRes)
	}
	if extOK.Result != "success" {
		t.Errorf("expected result=success, got %s", extOK.Result)
	}
	if extOK.Message.ID != msgID {
		t.Errorf("expected message ID=%s, got %s", msgID, extOK.Message.ID)
	}

	// Extend timeout for nonexistent message
	extRes2, err := client.ExtendMessageTimeout(ctx, message.ExtendMessageTimeoutParams{
		QueueName: message.QueueName(queueName),
		MessageId: nonexistentUUID,
	})
	if err != nil {
		t.Fatalf("extend timeout request failed: %v", err)
	}
	if _, ok := extRes2.(*message.ExtendMessageTimeoutNotFound); !ok {
		t.Fatalf("expected ExtendMessageTimeoutNotFound, got %T", extRes2)
	}
}

func TestVisibilityTimeout(t *testing.T) {
	srv := localserver.NewServer()
	defer srv.Close()

	client := newTestClient(t, srv.URL(), "test-api-key")
	ctx := context.Background()
	queueName := "test-queue"

	// Send a message
	_, err := client.SendMessage(ctx, &message.SendRequest{Content: message.MessageContent(b64("visibility test"))}, message.SendMessageParams{QueueName: message.QueueName(queueName)})
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}

	// Receive the message (makes it invisible)
	recvRes, err := client.ReceiveMessage(ctx, message.ReceiveMessageParams{QueueName: message.QueueName(queueName)})
	if err != nil {
		t.Fatalf("receive failed: %v", err)
	}
	recvOK := recvRes.(*message.ReceiveMessageOK)
	if len(recvOK.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(recvOK.Messages))
	}

	// Second receive should return empty (message is invisible)
	recvRes2, err := client.ReceiveMessage(ctx, message.ReceiveMessageParams{QueueName: message.QueueName(queueName)})
	if err != nil {
		t.Fatalf("second receive failed: %v", err)
	}
	recvOK2 := recvRes2.(*message.ReceiveMessageOK)
	if len(recvOK2.Messages) != 0 {
		t.Errorf("expected 0 messages (invisible), got %d", len(recvOK2.Messages))
	}
}
