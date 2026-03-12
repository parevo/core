package postgres

import (
	"context"
	"database/sql"

	"github.com/parevo/core/storage"
)

type PostgresSessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *PostgresSessionStore {
	return &PostgresSessionStore{db: db}
}

func (s *PostgresSessionStore) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_sessions SET revoked = TRUE WHERE session_id = $1`, sessionID)
	return err
}

func (s *PostgresSessionStore) IsSessionRevoked(ctx context.Context, sessionID string) (bool, error) {
	var revoked bool
	err := s.db.QueryRowContext(ctx, `SELECT revoked FROM parevo_sessions WHERE session_id = $1`, sessionID).Scan(&revoked)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return revoked, err
}

func (s *PostgresSessionStore) BindSessionToUser(ctx context.Context, userID, sessionID string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parevo_sessions (session_id, user_id, revoked) VALUES ($1, $2, FALSE) ON CONFLICT (session_id) DO UPDATE SET user_id = $2`,
		sessionID, userID)
	return err
}

func (s *PostgresSessionStore) RevokeAllSessionsByUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_sessions SET revoked = TRUE WHERE user_id = $1`, userID)
	return err
}

var _ storage.UserSessionStore = (*PostgresSessionStore)(nil)
