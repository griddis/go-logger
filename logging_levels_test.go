package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger_Levels_Formats(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		level          string
		logFunc        func(l Logger, msg string)
		msg            string
		expectedSubstr []string
	}{
		{
			name:   "JSON Info",
			format: "json",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Info(msg)
			},
			msg:            "info message",
			expectedSubstr: []string{`"level":"info"`, `"message":"info message"`},
		},
		{
			name:   "JSON Warn",
			format: "json",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Warn(msg)
			},
			msg:            "warn message",
			expectedSubstr: []string{`"level":"warning"`, `"message":"warn message"`},
		},
		{
			name:   "JSON Error",
			format: "json",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Error(msg)
			},
			msg:            "error message",
			expectedSubstr: []string{`"level":"error"`, `"message":"error message"`},
		},
		{
			name:   "JSON Debug",
			format: "json",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Debug(msg)
			},
			msg:            "debug message",
			expectedSubstr: []string{`"level":"debug"`, `"message":"debug message"`},
		},
		{
			name:   "Plaintext Info",
			format: "plaintext",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Info(msg)
			},
			msg:            "info message",
			expectedSubstr: []string{"✔", "info message"},
		},
		{
			name:   "Plaintext Warn",
			format: "plaintext",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Warn(msg)
			},
			msg:            "warn message",
			expectedSubstr: []string{"⚠", "warn message"},
		},
		{
			name:   "Plaintext Error",
			format: "plaintext",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Error(msg)
			},
			msg:            "error message",
			expectedSubstr: []string{"✘", "error message"},
		},
		{
			name:   "Plaintext Debug",
			format: "plaintext",
			level:  "debug",
			logFunc: func(l Logger, msg string) {
				l.Debug(msg)
			},
			msg:            "debug message",
			expectedSubstr: []string{"•", "debug message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := &Config{
				Level:  tt.level,
				Format: tt.format,
				Writer: &buf,
			}
			l := NewLogger(cfg)
			tt.logFunc(l, tt.msg)

			output := buf.String()
			for _, substr := range tt.expectedSubstr {
				if !strings.Contains(output, substr) {
					t.Errorf("expected output %q to contain %q", output, substr)
				}
			}
		})
	}
}
