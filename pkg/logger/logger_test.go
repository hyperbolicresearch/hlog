package logger

import (
	"bytes"
	"testing"

	"github.com/hyperbolicresearch/hlog/internal/logs"
)

func TestLogger(t *testing.T) {
	var buff bytes.Buffer
	logger := New(DEBUG, &buff)
	// string
	dataStr := "this is literally a test string"
	// bytes
	dataBytes := []byte(dataStr)
	// core.Log
	log := logs.Log{
		Channel:   "channel",
		LogId:     "0000-0000-0000-0000",
		SenderId:  "sender_id",
		Timestamp: 1618304400,
		Level:     "debug",
		Message:   "this is a test message",
		Data:      map[string]interface{}{"foo": "bar"},
	}

	logger.Log(dataStr)
	logger.Log(dataBytes)
	logger.Log(log)

	tests := []struct {
		name     string
		expected string
		input    interface{}
	}{
		{name: "Logging string", expected: "foo\n", input: "foo"},
		{name: "Logging bytes", expected: "foo\n", input: []byte("foo")},
		{
			name:     "Logging core.Log",
			expected: "[channel | sender_id] 2021-04-13 12:00:00 DEBUG " + "\033[37mthis is a test message " + `{"foo":"bar"}` + "\033[0m\n",
			input:    log,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff.Reset()
			logger.Log(tt.input)
			output := buff.String()
			if tt.expected != output {
				t.Errorf("Expected=%s, Got=%s", tt.expected, output)
			}
		})
	}
}
