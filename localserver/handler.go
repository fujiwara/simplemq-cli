package localserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	queueNamePattern      = regexp.MustCompile(`^[0-9a-zA-Z]+(-[0-9a-zA-Z]+)*$`)
	messageContentPattern = regexp.MustCompile(`^[0-9a-zA-Z+/=]*$`)
	messageIDPattern      = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

const (
	queueNameMinLength   = 5
	queueNameMaxLength   = 64
	messageContentMaxLen = 256000
)

func validateQueueName(name string) error {
	if len(name) < queueNameMinLength || len(name) > queueNameMaxLength {
		return fmt.Errorf("queue name must be between %d and %d characters", queueNameMinLength, queueNameMaxLength)
	}
	if !queueNamePattern.MatchString(name) {
		return fmt.Errorf("queue name must match pattern %s", queueNamePattern.String())
	}
	return nil
}

func validateMessageContent(content string) error {
	if len(content) > messageContentMaxLen {
		return fmt.Errorf("message content must not exceed %d characters", messageContentMaxLen)
	}
	if !messageContentPattern.MatchString(content) {
		return fmt.Errorf("message content must be base64 encoded")
	}
	return nil
}

func validateMessageID(id string) error {
	if !messageIDPattern.MatchString(id) {
		return fmt.Errorf("invalid message ID format")
	}
	return nil
}

func (s *Server) buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/queues/{queueName}/messages", s.authMiddleware(s.handleSend))
	mux.HandleFunc("GET /v1/queues/{queueName}/messages", s.authMiddleware(s.handleReceive))
	mux.HandleFunc("PUT /v1/queues/{queueName}/messages/{messageID}", s.authMiddleware(s.handleExtendTimeout))
	mux.HandleFunc("DELETE /v1/queues/{queueName}/messages/{messageID}", s.authMiddleware(s.handleDelete))
	return mux
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") == "" {
			writeError(w, http.StatusUnauthorized, "authorization required")
			return
		}
		next(w, r)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
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

type sendMessageResponse struct {
	Result  string             `json:"result"`
	Message newMessageResponse `json:"message"`
}

type receiveMessagesResponse struct {
	Result   string            `json:"result"`
	Messages []messageResponse `json:"messages"`
}

type singleMessageResponse struct {
	Result  string          `json:"result"`
	Message messageResponse `json:"message"`
}

type successResponse struct {
	Result string `json:"result"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queueName")
	if err := validateQueueName(queueName); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	q := s.store.getQueue(queueName)
	now := time.Now()

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
	if err := validateMessageContent(req.Content); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg := q.send(req.Content, now)
	writeJSON(w, http.StatusOK, sendMessageResponse{
		Result: "success",
		Message: newMessageResponse{
			ID:        msg.ID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.UnixMilli(),
			UpdatedAt: msg.UpdatedAt.UnixMilli(),
			ExpiresAt: msg.ExpiresAt.UnixMilli(),
		},
	})
}

func (s *Server) handleReceive(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queueName")
	if err := validateQueueName(queueName); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	q := s.store.getQueue(queueName)
	now := time.Now()

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
	writeJSON(w, http.StatusOK, receiveMessagesResponse{
		Result:   "success",
		Messages: messages,
	})
}

func (s *Server) handleExtendTimeout(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queueName")
	messageID := r.PathValue("messageID")
	if err := validateQueueName(queueName); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateMessageID(messageID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	q := s.store.getQueue(queueName)
	now := time.Now()

	msg, err := q.extendTimeout(messageID, now)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, singleMessageResponse{
		Result: "success",
		Message: messageResponse{
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

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queueName")
	messageID := r.PathValue("messageID")
	if err := validateQueueName(queueName); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateMessageID(messageID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	q := s.store.getQueue(queueName)

	if err := q.delete(messageID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, successResponse{
		Result: "success",
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
