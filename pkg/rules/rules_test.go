package rules

import "testing"

func TestCheckMessage(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		expected []string
	}{
		{
			name:     "valid message",
			msg:      "starting server",
			expected: nil,
		},
		{
			name: "uppercase start",
			msg:  "Starting server",
			expected: []string{
				"log message must start with a lowercase letter",
			},
		},
		{
			name: "non english",
			msg:  "запуск сервера",
			expected: []string{
				"log message must be in English only",
			},
		},
		{
			name: "special chars",
			msg:  "connection failed!!!",
			expected: []string{
				"log message must not contain special characters or emoji",
			},
		},
		{
			name: "sensitive data",
			msg:  "token: abc123",
			expected: []string{
				"log message must not contain sensitive data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckMessage(tt.msg)

			if len(got) != len(tt.expected) {
				t.Fatalf("unexpected issue count: got %d, want %d, issues=%v", len(got), len(tt.expected), got)
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Fatalf("unexpected issue at index %d: got %q, want %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
