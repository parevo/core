// Package export provides GDPR/compliance data export for user data portability.
//
// # Usage
//
//	payload := export.NewPayload(userID)
//	payload.Profile = map[string]any{"email": "...", "name": "..."}
//	payload.Sessions = sessionsFromStore(ctx, userID)
//	jsonBytes, _ := export.ToJSON(payload)
//	// veya blob'a yaz: export.ToBlob(ctx, store, bucket, key, payload)
package export

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/parevo/core/blob"
)

// Payload holds user data for export (GDPR Article 20).
type Payload struct {
	UserID      string         `json:"user_id"`
	Profile     map[string]any `json:"profile,omitempty"`
	Sessions    []map[string]any `json:"sessions,omitempty"`
	Consents    []map[string]any `json:"consents,omitempty"`
	Permissions []map[string]any `json:"permissions,omitempty"`
	ExportedAt  time.Time      `json:"exported_at"`
}

// NewPayload creates an export payload for the given user.
func NewPayload(userID string) *Payload {
	return &Payload{
		UserID:     userID,
		ExportedAt: time.Now().UTC(),
	}
}

// ToJSON serializes the payload to JSON.
func ToJSON(p *Payload) ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// ToReader returns an io.Reader for streaming.
func ToReader(p *Payload) (io.Reader, error) {
	b, err := ToJSON(p)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// ToBlob writes the export to blob storage.
func ToBlob(ctx context.Context, store blob.Store, bucket, key string, p *Payload) error {
	r, err := ToReader(p)
	if err != nil {
		return err
	}
	return store.Put(ctx, bucket, key, r, "application/json")
}
