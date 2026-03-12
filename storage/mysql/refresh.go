package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/parevo/core/storage"
)

// RefreshStore implements storage.UserRefreshTokenStore for MySQL.
type RefreshStore struct {
	db *sql.DB
}

// NewRefreshStore creates a MySQL RefreshStore.
func NewRefreshStore(db *sql.DB) *RefreshStore {
	return &RefreshStore{db: db}
}

// MarkIssued records a refresh token as issued.
func (s *RefreshStore) MarkIssued(ctx context.Context, sessionID, tokenID string, expiresAt time.Time) error {
	var userID string
	err := s.db.QueryRowContext(ctx, `SELECT user_id FROM parevo_sessions WHERE session_id = ?`, sessionID).Scan(&userID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO parevo_refresh_tokens (token_id, session_id, user_id, expires_at) VALUES (?, ?, ?, ?)`,
		tokenID, sessionID, userID, expiresAt)
	return err
}

// IsUsed checks if the token was already used.
func (s *RefreshStore) IsUsed(ctx context.Context, tokenID string) (bool, error) {
	var replacedBy sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT replaced_by FROM parevo_refresh_tokens WHERE token_id = ?`, tokenID).Scan(&replacedBy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return replacedBy.Valid, nil
}

// MarkUsed marks the token as used.
func (s *RefreshStore) MarkUsed(ctx context.Context, tokenID, replacedBy string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_refresh_tokens SET replaced_by = ? WHERE token_id = ?`, replacedBy, tokenID)
	return err
}

// RevokeSession revokes all refresh tokens for the session.
func (s *RefreshStore) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_refresh_tokens SET replaced_by = 'revoked' WHERE session_id = ?`, sessionID)
	return err
}

// BindSessionToUser is a no-op for MySQL (user_id comes from sessions).
func (s *RefreshStore) BindSessionToUser(_ context.Context, _, _ string) error {
	return nil
}

// RevokeAllByUser revokes all refresh tokens for the user.
func (s *RefreshStore) RevokeAllByUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_refresh_tokens SET replaced_by = 'revoked' WHERE user_id = ?`, userID)
	return err
}

var _ storage.UserRefreshTokenStore = (*RefreshStore)(nil)
