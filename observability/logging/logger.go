package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "info"
	}
}

type contextKey string

const (
	RequestIDContextKey contextKey = "parevo.request_id"
	TenantIDContextKey  contextKey = "parevo.tenant_id"
	UserIDContextKey    contextKey = "parevo.user_id"
)

type Config struct {
	Environment string
	Service     string
	MinLevel    Level
	Writer      io.Writer
	RedactKeys  []string
}

type Logger struct {
	cfg       Config
	redactSet map[string]struct{}
	mu        sync.Mutex
}

func New(cfg Config) *Logger {
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}
	if cfg.Service == "" {
		cfg.Service = "parevo"
	}
	redact := map[string]struct{}{
		"authorization": {},
		"token":         {},
		"access_token":  {},
		"refresh_token": {},
		"password":      {},
		"secret":        {},
	}
	for _, k := range cfg.RedactKeys {
		redact[strings.ToLower(strings.TrimSpace(k))] = struct{}{}
	}
	return &Logger{
		cfg:       cfg,
		redactSet: redact,
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, LevelDebug, msg, fields)
}
func (l *Logger) Info(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, LevelInfo, msg, fields)
}
func (l *Logger) Warn(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, LevelWarn, msg, fields)
}
func (l *Logger) Error(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, LevelError, msg, fields)
}

func (l *Logger) log(ctx context.Context, level Level, msg string, fields map[string]any) {
	if level < l.cfg.MinLevel {
		return
	}
	entry := map[string]any{
		"ts":          time.Now().UTC().Format(time.RFC3339Nano),
		"level":       level.String(),
		"service":     l.cfg.Service,
		"environment": l.cfg.Environment,
		"message":     msg,
	}
	for k, v := range fields {
		entry[k] = l.sanitize(k, v)
	}
	if ctx != nil {
		if requestID, ok := ctx.Value(RequestIDContextKey).(string); ok && requestID != "" {
			entry["request_id"] = requestID
		}
		if tenantID, ok := ctx.Value(TenantIDContextKey).(string); ok && tenantID != "" {
			entry["tenant_id"] = tenantID
		}
		if userID, ok := ctx.Value(UserIDContextKey).(string); ok && userID != "" {
			entry["user_id"] = userID
		}
	}
	l.write(entry)
}

func (l *Logger) sanitize(key string, value any) any {
	if _, ok := l.redactSet[strings.ToLower(strings.TrimSpace(key))]; ok {
		return "[REDACTED]"
	}
	return value
}

func (l *Logger) write(entry map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if strings.EqualFold(l.cfg.Environment, "production") {
		payload, _ := json.Marshal(entry)
		_, _ = fmt.Fprintln(l.cfg.Writer, string(payload))
		return
	}

	keys := make([]string, 0, len(entry))
	for k := range entry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, entry[k]))
	}
	_, _ = fmt.Fprintln(l.cfg.Writer, strings.Join(parts, " "))
}
