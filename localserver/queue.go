package localserver

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	defaultVisibilityTimeout = 30 * time.Second
	defaultMessageExpiration = 4 * 24 * time.Hour
)

type storedMessage struct {
	ID                  string
	Content             string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ExpiresAt           time.Time
	AcquiredAt          time.Time
	VisibilityTimeoutAt time.Time
}

type queue struct {
	mu                sync.Mutex
	messages          []*storedMessage
	visibilityTimeout time.Duration
	messageExpiration time.Duration
}

func newQueue() *queue {
	return &queue{
		visibilityTimeout: defaultVisibilityTimeout,
		messageExpiration: defaultMessageExpiration,
	}
}

func (q *queue) send(content string, now time.Time) storedMessage {
	q.mu.Lock()
	defer q.mu.Unlock()

	msg := &storedMessage{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(q.messageExpiration),
	}
	q.messages = append(q.messages, msg)
	return *msg
}

func (q *queue) receive(now time.Time) (storedMessage, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.compact(now)

	for _, msg := range q.messages {
		if !msg.VisibilityTimeoutAt.IsZero() && now.Before(msg.VisibilityTimeoutAt) {
			continue
		}
		msg.AcquiredAt = now
		msg.UpdatedAt = now
		msg.VisibilityTimeoutAt = now.Add(q.visibilityTimeout)
		return *msg, true
	}
	return storedMessage{}, false
}

// compact removes expired messages from the slice. Must be called with mu held.
func (q *queue) compact(now time.Time) {
	n := 0
	for _, msg := range q.messages {
		if now.Before(msg.ExpiresAt) {
			q.messages[n] = msg
			n++
		}
	}
	clear(q.messages[n:])
	q.messages = q.messages[:n]
}

func (q *queue) extendTimeout(id string, now time.Time) (storedMessage, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, msg := range q.messages {
		if msg.ID == id {
			msg.VisibilityTimeoutAt = now.Add(q.visibilityTimeout)
			msg.UpdatedAt = now
			return *msg, nil
		}
	}
	return storedMessage{}, fmt.Errorf("message not found: %s", id)
}

func (q *queue) delete(id string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, msg := range q.messages {
		if msg.ID == id {
			q.messages = append(q.messages[:i], q.messages[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("message not found: %s", id)
}

// Store manages named queues. Queues are created on first access.
type Store struct {
	mu     sync.Mutex
	queues map[string]*queue
}

// NewStore creates a new empty Store.
func NewStore() *Store {
	return &Store{
		queues: make(map[string]*queue),
	}
}

func (s *Store) getQueue(name string) *queue {
	s.mu.Lock()
	defer s.mu.Unlock()

	q, ok := s.queues[name]
	if !ok {
		q = newQueue()
		s.queues[name] = q
	}
	return q
}
