package localserver

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var pathPattern = regexp.MustCompile(`^/v1/queues/([^/]+)/messages(?:/([^/]+))?$`)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Auth check: Bearer token must be present
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") == "" {
		writeError(w, http.StatusUnauthorized, "authorization required")
		return
	}

	matches := pathPattern.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	queueName := matches[1]
	messageID := matches[2]
	q := s.store.getQueue(queueName)
	now := time.Now()

	switch r.Method {
	case http.MethodPost:
		if messageID != "" {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		s.handleSend(w, r, q, now)
	case http.MethodGet:
		if messageID != "" {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		s.handleReceive(w, q, now)
	case http.MethodPut:
		if messageID == "" {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		s.handleExtendTimeout(w, q, messageID, now)
	case http.MethodDelete:
		if messageID == "" {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		s.handleDelete(w, q, messageID)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

type sendRequest struct {
	Content string `json:"content"`
}

type newMessageResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	ExpiresAt int64  `json:"expires_at"`
}

type messageResponse struct {
	ID                  string `json:"id"`
	Content             string `json:"content"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
	ExpiresAt           int64  `json:"expires_at"`
	AcquiredAt          int64  `json:"acquired_at"`
	VisibilityTimeoutAt int64  `json:"visibility_timeout_at"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request, q *queue, now time.Time) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read body")
		return
	}
	var req sendRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	msg := q.send(req.Content, now)
	writeJSON(w, http.StatusOK, map[string]any{
		"result": "success",
		"message": newMessageResponse{
			ID:        msg.ID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.UnixMilli(),
			UpdatedAt: msg.UpdatedAt.UnixMilli(),
			ExpiresAt: msg.ExpiresAt.UnixMilli(),
		},
	})
}

func (s *Server) handleReceive(w http.ResponseWriter, q *queue, now time.Time) {
	msg := q.receive(now)
	messages := []messageResponse{}
	if msg != nil {
		messages = append(messages, messageResponse{
			ID:                  msg.ID,
			Content:             msg.Content,
			CreatedAt:           msg.CreatedAt.UnixMilli(),
			UpdatedAt:           msg.UpdatedAt.UnixMilli(),
			ExpiresAt:           msg.ExpiresAt.UnixMilli(),
			AcquiredAt:          msg.AcquiredAt.UnixMilli(),
			VisibilityTimeoutAt: msg.VisibilityTimeoutAt.UnixMilli(),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"result":   "success",
		"messages": messages,
	})
}

func (s *Server) handleExtendTimeout(w http.ResponseWriter, q *queue, messageID string, now time.Time) {
	msg, err := q.extendTimeout(messageID, now)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"result": "success",
		"message": messageResponse{
			ID:                  msg.ID,
			Content:             msg.Content,
			CreatedAt:           msg.CreatedAt.UnixMilli(),
			UpdatedAt:           msg.UpdatedAt.UnixMilli(),
			ExpiresAt:           msg.ExpiresAt.UnixMilli(),
			AcquiredAt:          msg.AcquiredAt.UnixMilli(),
			VisibilityTimeoutAt: msg.VisibilityTimeoutAt.UnixMilli(),
		},
	})
}

func (s *Server) handleDelete(w http.ResponseWriter, q *queue, messageID string) {
	if err := q.delete(messageID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"result": "success",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{
		Code:    status,
		Message: msg,
	})
}
