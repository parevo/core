package logging

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLoggerRedaction(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Config{
		Environment: "production",
		Service:     "parevo-test",
		Writer:      &buf,
		MinLevel:    LevelDebug,
	})
	ctx := context.WithValue(context.Background(), RequestIDContextKey, "req-1")
	logger.Info(ctx, "login", map[string]any{
		"token": "very-secret",
		"email": "u@example.com",
	})

	out := buf.String()
	if !strings.Contains(out, `"request_id":"req-1"`) {
		t.Fatalf("expected request id in log")
	}
	if strings.Contains(out, "very-secret") {
		t.Fatalf("expected token to be redacted")
	}
}
