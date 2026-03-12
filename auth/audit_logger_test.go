package auth

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/parevo/core/observability/logging"
)

func TestStructuredAuditLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(logging.Config{
		Environment: "production",
		Service:     "parevo-auth-test",
		Writer:      &buf,
		MinLevel:    logging.LevelDebug,
	})
	audit := NewStructuredAuditLogger(logger)
	ctx := WithRequestID(context.Background(), "req-42")
	audit.Log(ctx, AuditAuthFailed, map[string]string{
		"reason": "invalid_token",
		"token":  "secret-token",
	})

	out := buf.String()
	if !strings.Contains(out, `"event":"auth_failed"`) {
		t.Fatalf("expected audit event in output")
	}
	if strings.Contains(out, "secret-token") {
		t.Fatalf("expected redaction for token")
	}
}
