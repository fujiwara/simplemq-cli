package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestSendMessageCommandValidate(t *testing.T) {
	tests := []struct {
		name    string
		cmd     SendMessageCommand
		wantErr string
	}{
		{
			name:    "content arg only",
			cmd:     SendMessageCommand{Content: "hello"},
			wantErr: "",
		},
		{
			name:    "stdin only",
			cmd:     SendMessageCommand{Stdin: true},
			wantErr: "",
		},
		{
			name:    "stdin with each-line",
			cmd:     SendMessageCommand{Stdin: true, EachLine: true},
			wantErr: "",
		},
		{
			name:    "stdin with each-json",
			cmd:     SendMessageCommand{Stdin: true, EachJSON: true},
			wantErr: "",
		},
		{
			name:    "file only",
			cmd:     SendMessageCommand{File: "test.txt"},
			wantErr: "",
		},
		{
			name:    "file with each-line",
			cmd:     SendMessageCommand{File: "test.txt", EachLine: true},
			wantErr: "",
		},
		{
			name:    "file with each-json",
			cmd:     SendMessageCommand{File: "test.txt", EachJSON: true},
			wantErr: "",
		},
		{
			name:    "no content no stdin no file",
			cmd:     SendMessageCommand{},
			wantErr: "<content> argument, --stdin, or --file is required",
		},
		{
			name:    "stdin and content",
			cmd:     SendMessageCommand{Stdin: true, Content: "hello"},
			wantErr: "--stdin/--file and <content> argument are mutually exclusive",
		},
		{
			name:    "file and content",
			cmd:     SendMessageCommand{File: "test.txt", Content: "hello"},
			wantErr: "--stdin/--file and <content> argument are mutually exclusive",
		},
		{
			name:    "each-line without stdin or file",
			cmd:     SendMessageCommand{EachLine: true, Content: "hello"},
			wantErr: "--each-line requires --stdin or --file",
		},
		{
			name:    "each-json without stdin or file",
			cmd:     SendMessageCommand{EachJSON: true, Content: "hello"},
			wantErr: "--each-json requires --stdin or --file",
		},
		{
			name:    "each-line and each-json",
			cmd:     SendMessageCommand{Stdin: true, EachLine: true, EachJSON: true},
			wantErr: "--each-line and --each-json are mutually exclusive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}

func TestSendEachLine(t *testing.T) {
	input := "line1\n\nline2\nline3\n"
	var sent []string
	sendOne := func(data []byte) error {
		sent = append(sent, string(data))
		return nil
	}
	if err := sendEachLine(strings.NewReader(input), sendOne); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"line1", "line2", "line3"}
	if len(sent) != len(expected) {
		t.Fatalf("expected %d messages, got %d", len(expected), len(sent))
	}
	for i, s := range sent {
		if s != expected[i] {
			t.Errorf("message[%d] = %q, want %q", i, s, expected[i])
		}
	}
}

func TestSendEachJSON(t *testing.T) {
	input := `{"key":"value1"} {"key":"value2"}` + "\n" + `[1,2,3]`
	var sent []string
	sendOne := func(data []byte) error {
		sent = append(sent, string(data))
		return nil
	}
	if err := sendEachJSON(strings.NewReader(input), sendOne); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{`{"key":"value1"}`, `{"key":"value2"}`, `[1,2,3]`}
	if len(sent) != len(expected) {
		t.Fatalf("expected %d messages, got %d: %v", len(expected), len(sent), sent)
	}
	for i, s := range sent {
		if s != expected[i] {
			t.Errorf("message[%d] = %q, want %q", i, s, expected[i])
		}
	}
}

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
