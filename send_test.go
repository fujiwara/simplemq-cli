package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		name     string
		input    io.Reader
		expected []byte
		wantErr  bool
	}{
		{
			name:     "simple string",
			input:    strings.NewReader("Hello, World!"),
			expected: []byte("Hello, World!"),
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    strings.NewReader(""),
			expected: []byte{},
			wantErr:  false,
		},
		{
			name:     "multiline input",
			input:    strings.NewReader("line1\nline2\nline3"),
			expected: []byte("line1\nline2\nline3"),
			wantErr:  false,
		},
		{
			name:     "binary data",
			input:    bytes.NewReader([]byte{0x00, 0x01, 0x02, 0xff, 0xfe}),
			expected: []byte{0x00, 0x01, 0x02, 0xff, 0xfe},
			wantErr:  false,
		},
		{
			name:     "Japanese text",
			input:    strings.NewReader("こんにちは世界"),
			expected: []byte("こんにちは世界"),
			wantErr:  false,
		},
		{
			name:     "large input",
			input:    strings.NewReader(strings.Repeat("a", 10000)),
			expected: []byte(strings.Repeat("a", 10000)),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("readInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("readInput() = %v, want %v", got, tt.expected)
			}
		})
	}
}
