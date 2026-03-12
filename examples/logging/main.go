package main

import (
	"context"
	"os"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/observability/logging"
)

func main() {
	env := os.Getenv("PAREVO_ENV")
	if env == "" {
		env = "development"
	}
	logger := logging.New(logging.Config{
		Environment: env,
		Service:     "parevo-example",
		MinLevel:    logging.LevelDebug,
	})
	audit := auth.NewStructuredAuditLogger(logger)

	ctx := auth.WithRequestID(context.Background(), "req-demo-1")
	audit.Log(ctx, auth.AuditAuthSuccess, map[string]string{
		"user_id":      "u1",
		"tenant_id":    "t1",
		"access_token": "should-be-redacted",
	})
}
