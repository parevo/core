package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/parevo/core/storage"
)

// SocialAccountStore implements storage.SocialAccountStore for Postgres.
type SocialAccountStore struct {
	db *sql.DB
}

// NewSocialAccountStore creates a Postgres SocialAccountStore.
func NewSocialAccountStore(db *sql.DB) *SocialAccountStore {
	return &SocialAccountStore{db: db}
}

// FindUserBySocial returns the userID for the given provider and provider user ID.
func (s *SocialAccountStore) FindUserBySocial(ctx context.Context, provider, providerUserID string) (userID string, found bool, err error) {
	err = s.db.QueryRowContext(ctx,
		`SELECT user_id FROM parevo_social_accounts WHERE provider = $1 AND provider_user_id = $2`,
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
	userID := "user:" + email
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO parevo_users (user_id, email, display_name) VALUES ($1, $2, $3)
		 ON CONFLICT (email) DO UPDATE SET display_name = EXCLUDED.display_name
		 RETURNING user_id`,
		userID, email, displayName).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("social: find or create user: %w", err)
	}
	return userID, nil
}

// LinkSocialAccount links a social identity to a user.
func (s *SocialAccountStore) LinkSocialAccount(ctx context.Context, userID string, identity storage.SocialIdentity) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parevo_social_accounts (provider, provider_user_id, user_id, email, email_verified, name, avatar_url)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (provider, provider_user_id) DO UPDATE SET user_id = $3, email = $4, email_verified = $5, name = $6, avatar_url = $7`,
		identity.Provider, identity.ProviderUserID, userID, identity.Email, identity.EmailVerified, identity.Name, identity.AvatarURL)
	return err
}

var _ storage.SocialAccountStore = (*SocialAccountStore)(nil)
