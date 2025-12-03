package cli

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/sacloud/simplemq-api-go/apis/v1/message"
)

func TestUnixToTime(t *testing.T) {
	tests := []struct {
		name     string
		msec     int64
		expected time.Time
	}{
		{
			name:     "zero",
			msec:     0,
			expected: time.Unix(0, 0),
		},
		{
			name:     "1 second",
			msec:     1000,
			expected: time.Unix(1, 0),
		},
		{
			name:     "with milliseconds",
			msec:     1234567,
			expected: time.Unix(1234, 567_000_000),
		},
		{
			name:     "real timestamp",
			msec:     1733140964820,
			expected: time.Unix(1733140964, 820_000_000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnixToTime(tt.msec)
			if !got.Equal(tt.expected) {
				t.Errorf("UnixToTime(%d) = %v, want %v", tt.msec, got, tt.expected)
			}
		})
	}
}

func TestConvertMessageContent(t *testing.T) {
	tests := []struct {
		name            string
		msg             message.Message
		expectedContent string
	}{
		{
			name: "valid base64 content",
			msg: message.Message{
				ID:                  "test-id-1",
				Content:             message.MessageContent(base64.StdEncoding.EncodeToString([]byte("Hello, World!"))),
				CreatedAt:           1733140964820,
				UpdatedAt:           1733140965000,
				ExpiresAt:           1733227364820,
				AcquiredAt:          1733140965100,
				VisibilityTimeoutAt: 1733140994820,
			},
			expectedContent: "Hello, World!",
		},
		{
			name: "valid base64 content with Japanese",
			msg: message.Message{
				ID:                  "test-id-2",
				Content:             message.MessageContent(base64.StdEncoding.EncodeToString([]byte("こんにちは世界"))),
				CreatedAt:           1733140964820,
				UpdatedAt:           1733140965000,
				ExpiresAt:           1733227364820,
				AcquiredAt:          1733140965100,
				VisibilityTimeoutAt: 1733140994820,
			},
			expectedContent: "こんにちは世界",
		},
		{
			name: "invalid base64 content returns as-is",
			msg: message.Message{
				ID:                  "test-id-3",
				Content:             "not-valid-base64!!!",
				CreatedAt:           1733140964820,
				UpdatedAt:           1733140965000,
				ExpiresAt:           1733227364820,
				AcquiredAt:          1733140965100,
				VisibilityTimeoutAt: 1733140994820,
			},
			expectedContent: "not-valid-base64!!!",
		},
		{
			name: "empty content",
			msg: message.Message{
				ID:                  "test-id-4",
				Content:             "",
				CreatedAt:           1733140964820,
				UpdatedAt:           1733140965000,
				ExpiresAt:           1733227364820,
				AcquiredAt:          1733140965100,
				VisibilityTimeoutAt: 1733140994820,
			},
			expectedContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertMessageContent(tt.msg)

			if got.ID != string(tt.msg.ID) {
				t.Errorf("ID = %v, want %v", got.ID, tt.msg.ID)
			}
			if got.Content != tt.expectedContent {
				t.Errorf("Content = %v, want %v", got.Content, tt.expectedContent)
			}
			if !got.CreatedAt.Equal(UnixToTime(tt.msg.CreatedAt)) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, UnixToTime(tt.msg.CreatedAt))
			}
			if !got.UpdatedAt.Equal(UnixToTime(tt.msg.UpdatedAt)) {
				t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, UnixToTime(tt.msg.UpdatedAt))
			}
			if !got.ExpiresAt.Equal(UnixToTime(tt.msg.ExpiresAt)) {
				t.Errorf("ExpiresAt = %v, want %v", got.ExpiresAt, UnixToTime(tt.msg.ExpiresAt))
			}
			if !got.AcquiredAt.Equal(UnixToTime(tt.msg.AcquiredAt)) {
				t.Errorf("AcquiredAt = %v, want %v", got.AcquiredAt, UnixToTime(tt.msg.AcquiredAt))
			}
			if !got.VisibilityTimeoutAt.Equal(UnixToTime(tt.msg.VisibilityTimeoutAt)) {
				t.Errorf("VisibilityTimeoutAt = %v, want %v", got.VisibilityTimeoutAt, UnixToTime(tt.msg.VisibilityTimeoutAt))
			}
		})
	}
}

func TestErrNotFound(t *testing.T) {
	err := ErrNotFound{Message: "test error message"}
	if err.Error() != "test error message" {
		t.Errorf("Error() = %v, want %v", err.Error(), "test error message")
	}
}
