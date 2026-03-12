package mysql

import (
	"context"
	"database/sql"

	"github.com/parevo/core/storage"
)

// SessionStore implements storage.UserSessionStore for MySQL.
type SessionStore struct {
	db *sql.DB
}

// NewSessionStore creates a MySQL SessionStore.
func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

// RevokeSession marks the session as revoked.
func (s *SessionStore) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_sessions SET revoked = TRUE WHERE session_id = ?`, sessionID)
	return err
}

// IsSessionRevoked checks if the session is revoked.
func (s *SessionStore) IsSessionRevoked(ctx context.Context, sessionID string) (bool, error) {
	var revoked bool
	err := s.db.QueryRowContext(ctx, `SELECT revoked FROM parevo_sessions WHERE session_id = ?`, sessionID).Scan(&revoked)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return revoked, err
}

// BindSessionToUser binds a session to a user.
func (s *SessionStore) BindSessionToUser(ctx context.Context, userID, sessionID string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parevo_sessions (session_id, user_id, revoked) VALUES (?, ?, FALSE)
		 ON DUPLICATE KEY UPDATE user_id = VALUES(user_id)`,
		sessionID, userID)
	return err
}

// RevokeAllSessionsByUser revokes all sessions for the user.
func (s *SessionStore) RevokeAllSessionsByUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_sessions SET revoked = TRUE WHERE user_id = ?`, userID)
	return err
}

// ListSessionsByUser returns session IDs for the user (for admin UI).
func (s *SessionStore) ListSessionsByUser(ctx context.Context, userID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT session_id FROM parevo_sessions WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var sid string
		if err := rows.Scan(&sid); err != nil {
			return nil, err
		}
		out = append(out, sid)
	}
	return out, rows.Err()
}

var _ storage.UserSessionStore = (*SessionStore)(nil)
var _ storage.SessionListStore = (*SessionStore)(nil)
