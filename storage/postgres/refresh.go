package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/parevo/core/storage"
)

type PostgresRefreshStore struct {
	db *sql.DB
}

func NewRefreshStore(db *sql.DB) *PostgresRefreshStore {
	return &PostgresRefreshStore{db: db}
}

func (s *PostgresRefreshStore) MarkIssued(ctx context.Context, sessionID, tokenID string, expiresAt time.Time) error {
	var userID string
	err := s.db.QueryRowContext(ctx, `SELECT user_id FROM parevo_sessions WHERE session_id = $1`, sessionID).Scan(&userID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO parevo_refresh_tokens (token_id, session_id, user_id, expires_at) VALUES ($1, $2, $3, $4)`,
		tokenID, sessionID, userID, expiresAt)
	return err
}

func (s *PostgresRefreshStore) IsUsed(ctx context.Context, tokenID string) (bool, error) {
	var replacedBy sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT replaced_by FROM parevo_refresh_tokens WHERE token_id = $1`, tokenID).Scan(&replacedBy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return replacedBy.Valid, nil
}

func (s *PostgresRefreshStore) MarkUsed(ctx context.Context, tokenID, replacedBy string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_refresh_tokens SET replaced_by = $1 WHERE token_id = $2`, replacedBy, tokenID)
	return err
}

func (s *PostgresRefreshStore) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_refresh_tokens SET replaced_by = 'revoked' WHERE session_id = $1`, sessionID)
	return err
}

func (s *PostgresRefreshStore) BindSessionToUser(ctx context.Context, userID, sessionID string) error {
	return nil
}

func (s *PostgresRefreshStore) RevokeAllByUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_refresh_tokens SET replaced_by = 'revoked' WHERE user_id = $1`, userID)
	return err
}

var _ storage.UserRefreshTokenStore = (*PostgresRefreshStore)(nil)
