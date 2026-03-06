//go:build sqlite

package localserver

import (
	"path/filepath"
	"testing"
	"time"
)

func TestSQLiteStoreSendReceiveDelete(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewSQLiteStore(dbPath, 30*time.Second, time.Hour)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	defer store.Close()

	now := time.Now().Truncate(time.Millisecond)
	queueName := "test-queue"

	// Send
	msg, err := store.Send(queueName, "hello", now)
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if msg.ID == "" {
		t.Error("expected non-empty message ID")
	}
	if msg.Content != "hello" {
		t.Errorf("expected content=hello, got %s", msg.Content)
	}

	// Receive
	msg2, ok, err := store.Receive(queueName, now)
	if err != nil {
		t.Fatalf("receive failed: %v", err)
	}
	if !ok {
		t.Fatal("expected to receive a message")
	}
	if msg2.ID != msg.ID {
		t.Errorf("expected message ID=%s, got %s", msg.ID, msg2.ID)
	}
	if msg2.Content != "hello" {
		t.Errorf("expected content=hello, got %s", msg2.Content)
	}
	if msg2.AcquiredAt.IsZero() {
		t.Error("expected non-zero acquired_at")
	}

	// Second receive should return nothing (visibility timeout)
	_, ok, err = store.Receive(queueName, now)
	if err != nil {
		t.Fatalf("second receive failed: %v", err)
	}
	if ok {
		t.Error("expected no message within visibility timeout")
	}

	// Delete
	if err := store.Delete(queueName, msg.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Delete nonexistent
	if err := store.Delete(queueName, "00000000-0000-0000-0000-000000000000"); err == nil {
		t.Error("expected error for nonexistent message")
	}
}

func TestSQLiteStoreVisibilityTimeoutRevisibility(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewSQLiteStore(dbPath, 2*time.Second, time.Hour)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	defer store.Close()

	now := time.Now().Truncate(time.Millisecond)
	queueName := "test-queue"

	store.Send(queueName, "hello", now)

	msg, ok, _ := store.Receive(queueName, now)
	if !ok {
		t.Fatal("expected to receive a message")
	}

	// Within timeout
	_, ok, _ = store.Receive(queueName, now.Add(1*time.Second))
	if ok {
		t.Error("expected no message within visibility timeout")
	}

	// After timeout
	msg2, ok, _ := store.Receive(queueName, now.Add(3*time.Second))
	if !ok {
		t.Fatal("expected message to be visible again after timeout")
	}
	if msg2.ID != msg.ID {
		t.Errorf("expected same message ID %s, got %s", msg.ID, msg2.ID)
	}
}

func TestSQLiteStoreMessageExpiration(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewSQLiteStore(dbPath, time.Second, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	defer store.Close()

	now := time.Now().Truncate(time.Millisecond)
	queueName := "test-queue"

	store.Send(queueName, "expires soon", now)

	// Before expiration
	_, ok, _ := store.Receive(queueName, now.Add(4*time.Second))
	if !ok {
		t.Fatal("expected to receive message before expiration")
	}

	// After expiration (and past visibility timeout)
	_, ok, _ = store.Receive(queueName, now.Add(6*time.Second))
	if ok {
		t.Error("expected message to be expired and removed")
	}
}

func TestSQLiteStoreExtendTimeout(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewSQLiteStore(dbPath, 30*time.Second, time.Hour)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	defer store.Close()

	now := time.Now().Truncate(time.Millisecond)
	queueName := "test-queue"

	msg, _ := store.Send(queueName, "hello", now)
	store.Receive(queueName, now)

	extended, err := store.ExtendTimeout(queueName, msg.ID, now.Add(10*time.Second))
	if err != nil {
		t.Fatalf("extend timeout failed: %v", err)
	}
	if extended.ID != msg.ID {
		t.Errorf("expected message ID=%s, got %s", msg.ID, extended.ID)
	}

	// Nonexistent
	_, err = store.ExtendTimeout(queueName, "00000000-0000-0000-0000-000000000000", now)
	if err == nil {
		t.Error("expected error for nonexistent message")
	}
}

func TestSQLiteStoreMultipleQueues(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewSQLiteStore(dbPath, 30*time.Second, time.Hour)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	defer store.Close()

	now := time.Now().Truncate(time.Millisecond)

	store.Send("queue-a", "msg-a", now)
	store.Send("queue-b", "msg-b", now)

	msgA, ok, _ := store.Receive("queue-a", now)
	if !ok || msgA.Content != "msg-a" {
		t.Errorf("expected msg-a from queue-a, got %v", msgA.Content)
	}

	msgB, ok, _ := store.Receive("queue-b", now)
	if !ok || msgB.Content != "msg-b" {
		t.Errorf("expected msg-b from queue-b, got %v", msgB.Content)
	}

	// queue-a should be empty now (visibility timeout)
	_, ok, _ = store.Receive("queue-a", now)
	if ok {
		t.Error("expected no more messages in queue-a")
	}
}

func TestSQLiteStorePersistence(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	now := time.Now().Truncate(time.Millisecond)

	// Write with first store instance
	store1, err := NewSQLiteStore(dbPath, 30*time.Second, time.Hour)
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	msg, _ := store1.Send("test-queue", "persistent", now)
	store1.Close()

	// Read with second store instance
	store2, err := NewSQLiteStore(dbPath, 30*time.Second, time.Hour)
	if err != nil {
		t.Fatalf("failed to reopen sqlite store: %v", err)
	}
	defer store2.Close()

	msg2, ok, err := store2.Receive("test-queue", now)
	if err != nil {
		t.Fatalf("receive failed: %v", err)
	}
	if !ok {
		t.Fatal("expected to receive persisted message")
	}
	if msg2.ID != msg.ID {
		t.Errorf("expected message ID=%s, got %s", msg.ID, msg2.ID)
	}
	if msg2.Content != "persistent" {
		t.Errorf("expected content=persistent, got %s", msg2.Content)
	}
}
