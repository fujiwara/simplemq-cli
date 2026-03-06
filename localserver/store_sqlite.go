//go:build sqlite

package localserver

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    queue_name TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    acquired_at INTEGER NOT NULL DEFAULT 0,
    visibility_timeout_at INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_messages_queue_name ON messages(queue_name);
`

// SQLiteStore is a Store backed by SQLite.
type SQLiteStore struct {
	db                *sql.DB
	visibilityTimeout time.Duration
	messageExpiration time.Duration
}

// NewSQLiteStore opens (or creates) the SQLite database at path and returns a Store.
func NewSQLiteStore(path string, visibilityTimeout, messageExpiration time.Duration) (*SQLiteStore, error) {
	if visibilityTimeout <= 0 {
		visibilityTimeout = defaultVisibilityTimeout
	}
	if messageExpiration <= 0 {
		messageExpiration = defaultMessageExpiration
	}

	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}
	// SQLite does not support concurrent writes well; serialize access.
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	slog.Info("sqlite store opened", "path", path)
	return &SQLiteStore{
		db:                db,
		visibilityTimeout: visibilityTimeout,
		messageExpiration: messageExpiration,
	}, nil
}

func (s *SQLiteStore) Send(queueName, content string, now time.Time) (storedMessage, error) {
	msg := storedMessage{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(s.messageExpiration),
	}
	query := `INSERT INTO messages (id, queue_name, content, created_at, updated_at, expires_at) VALUES (?, ?, ?, ?, ?, ?)`
	args := []any{msg.ID, queueName, msg.Content, msg.CreatedAt.UnixMilli(), msg.UpdatedAt.UnixMilli(), msg.ExpiresAt.UnixMilli()}
	slog.Debug("sqlite exec", "sql", query, "args", args)
	_, err := s.db.Exec(query, args...)
	if err != nil {
		return storedMessage{}, fmt.Errorf("failed to insert message: %w", err)
	}
	return msg, nil
}

func (s *SQLiteStore) Receive(queueName string, now time.Time) (storedMessage, bool, error) {
	// Delete expired messages first
	compactQuery := `DELETE FROM messages WHERE queue_name = ? AND expires_at <= ?`
	compactArgs := []any{queueName, now.UnixMilli()}
	slog.Debug("sqlite exec", "sql", compactQuery, "args", compactArgs)
	if _, err := s.db.Exec(compactQuery, compactArgs...); err != nil {
		return storedMessage{}, false, fmt.Errorf("failed to compact expired messages: %w", err)
	}

	nowMilli := now.UnixMilli()
	visibilityTimeoutAtMilli := now.Add(s.visibilityTimeout).UnixMilli()

	// Find and atomically update the first available message
	query := `UPDATE messages SET acquired_at = ?, updated_at = ?, visibility_timeout_at = ?
		 WHERE rowid = (
		     SELECT rowid FROM messages
		     WHERE queue_name = ? AND expires_at > ? AND (visibility_timeout_at = 0 OR visibility_timeout_at <= ?)
		     ORDER BY rowid ASC LIMIT 1
		 )
		 RETURNING id, content, created_at, updated_at, expires_at, acquired_at, visibility_timeout_at`
	args := []any{nowMilli, nowMilli, visibilityTimeoutAtMilli, queueName, nowMilli, nowMilli}
	slog.Debug("sqlite query", "sql", query, "args", args)
	row := s.db.QueryRow(query, args...)

	var msg storedMessage
	var createdAt, updatedAt, expiresAt, acquiredAt, visibilityTimeoutAt int64
	err := row.Scan(&msg.ID, &msg.Content, &createdAt, &updatedAt, &expiresAt, &acquiredAt, &visibilityTimeoutAt)
	if err == sql.ErrNoRows {
		return storedMessage{}, false, nil
	}
	if err != nil {
		return storedMessage{}, false, fmt.Errorf("failed to receive message: %w", err)
	}
	msg.CreatedAt = time.UnixMilli(createdAt)
	msg.UpdatedAt = time.UnixMilli(updatedAt)
	msg.ExpiresAt = time.UnixMilli(expiresAt)
	msg.AcquiredAt = time.UnixMilli(acquiredAt)
	msg.VisibilityTimeoutAt = time.UnixMilli(visibilityTimeoutAt)
	return msg, true, nil
}

func (s *SQLiteStore) ExtendTimeout(queueName, id string, now time.Time) (storedMessage, error) {
	nowMilli := now.UnixMilli()
	visibilityTimeoutAtMilli := now.Add(s.visibilityTimeout).UnixMilli()

	query := `UPDATE messages SET visibility_timeout_at = ?, updated_at = ?
		 WHERE queue_name = ? AND id = ?
		 RETURNING id, content, created_at, updated_at, expires_at, acquired_at, visibility_timeout_at`
	args := []any{visibilityTimeoutAtMilli, nowMilli, queueName, id}
	slog.Debug("sqlite query", "sql", query, "args", args)
	row := s.db.QueryRow(query, args...)

	var msg storedMessage
	var createdAt, updatedAt, expiresAt, acquiredAt, visibilityTimeoutAt int64
	err := row.Scan(&msg.ID, &msg.Content, &createdAt, &updatedAt, &expiresAt, &acquiredAt, &visibilityTimeoutAt)
	if err == sql.ErrNoRows {
		return storedMessage{}, fmt.Errorf("message not found: %s", id)
	}
	if err != nil {
		return storedMessage{}, fmt.Errorf("failed to extend timeout: %w", err)
	}
	msg.CreatedAt = time.UnixMilli(createdAt)
	msg.UpdatedAt = time.UnixMilli(updatedAt)
	msg.ExpiresAt = time.UnixMilli(expiresAt)
	msg.AcquiredAt = time.UnixMilli(acquiredAt)
	msg.VisibilityTimeoutAt = time.UnixMilli(visibilityTimeoutAt)
	return msg, nil
}

func (s *SQLiteStore) Delete(queueName, id string) error {
	query := `DELETE FROM messages WHERE queue_name = ? AND id = ?`
	args := []any{queueName, id}
	slog.Debug("sqlite exec", "sql", query, "args", args)
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("message not found: %s", id)
	}
	return nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
