package postgres

import (
	"context"
	"database/sql"

	"github.com/parevo/core/storage"
)

// UserStore implements storage.UserListStore for Postgres.
type UserStore struct {
	db *sql.DB
}

// NewUserStore creates a Postgres UserStore.
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

// ListUsers returns users from parevo_users (for admin UI).
func (s *UserStore) ListUsers(ctx context.Context) ([]storage.UserInfo, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT user_id, email, display_name, created_at FROM parevo_users`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []storage.UserInfo
	for rows.Next() {
		var u storage.UserInfo
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

var _ storage.UserListStore = (*UserStore)(nil)
