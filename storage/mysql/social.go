package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/parevo/core/storage"
)

// SocialAccountStore implements storage.SocialAccountStore for MySQL.
type SocialAccountStore struct {
	db *sql.DB
}

// NewSocialAccountStore creates a MySQL SocialAccountStore.
func NewSocialAccountStore(db *sql.DB) *SocialAccountStore {
	return &SocialAccountStore{db: db}
}

// FindUserBySocial returns the userID for the given provider and provider user ID.
func (s *SocialAccountStore) FindUserBySocial(ctx context.Context, provider, providerUserID string) (userID string, found bool, err error) {
	err = s.db.QueryRowContext(ctx,
		`SELECT user_id FROM parevo_social_accounts WHERE provider = ? AND provider_user_id = ?`,
		provider, providerUserID).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return userID, true, nil
}

// FindOrCreateUserByEmail finds or creates a user by email.
func (s *SocialAccountStore) FindOrCreateUserByEmail(ctx context.Context, email, displayName string) (string, error) {
	var userID string
	err := s.db.QueryRowContext(ctx, `SELECT user_id FROM parevo_users WHERE email = ?`, email).Scan(&userID)
	if err == nil {
		return userID, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}
	userID = "user:" + email
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO parevo_users (user_id, email, display_name) VALUES (?, ?, ?)`,
		userID, email, displayName)
	if err != nil {
		return "", fmt.Errorf("social: create user: %w", err)
	}
	return userID, nil
}

// LinkSocialAccount links a social identity to a user.
func (s *SocialAccountStore) LinkSocialAccount(ctx context.Context, userID string, identity storage.SocialIdentity) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parevo_social_accounts (provider, provider_user_id, user_id, email, email_verified, name, avatar_url)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE user_id = VALUES(user_id), email = VALUES(email), email_verified = VALUES(email_verified), name = VALUES(name), avatar_url = VALUES(avatar_url)`,
		identity.Provider, identity.ProviderUserID, userID, identity.Email, identity.EmailVerified, identity.Name, identity.AvatarURL)
	return err
}

var _ storage.SocialAccountStore = (*SocialAccountStore)(nil)
