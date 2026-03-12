package storage

import (
	"context"
	"time"
)

type Subject struct {
	ID    string
	Roles []string
}

type TenantScope struct {
	SubjectID string
	TenantID  string
}

type PermissionGrant struct {
	SubjectID  string
	TenantID   string
	Permission string
	IsGranted  bool
}

// TenantStore is implemented by project-specific DB adapters.
type TenantStore interface {
	ResolveSubjectTenants(ctx context.Context, subjectID string) ([]string, error)
}

type TenantStatus string

const (
	TenantStatusActive   TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusDeleted  TenantStatus = "deleted"
)

type TenantLifecycleStore interface {
	Create(ctx context.Context, tenantID, name string, ownerID string) error
	Suspend(ctx context.Context, tenantID string) error
	Resume(ctx context.Context, tenantID string) error
	Delete(ctx context.Context, tenantID string) error
	Status(ctx context.Context, tenantID string) (TenantStatus, error)
}

// PermissionStore is implemented by project-specific DB adapters.
type PermissionStore interface {
	HasPermission(ctx context.Context, subjectID, tenantID, permission string, roles []string) (bool, error)
}

// RoleHierarchyStore returns role parents for hierarchy. Optional.
type RoleHierarchyStore interface {
	ParentRoles(ctx context.Context, role string) ([]string, error)
}

type SessionStore interface {
	RevokeSession(ctx context.Context, sessionID string) error
	IsSessionRevoked(ctx context.Context, sessionID string) (bool, error)
}

type UserSessionStore interface {
	SessionStore
	BindSessionToUser(ctx context.Context, userID, sessionID string) error
	RevokeAllSessionsByUser(ctx context.Context, userID string) error
}

type RefreshTokenStore interface {
	MarkIssued(ctx context.Context, sessionID, tokenID string, expiresAt time.Time) error
	IsUsed(ctx context.Context, tokenID string) (bool, error)
	MarkUsed(ctx context.Context, tokenID, replacedBy string) error
	RevokeSession(ctx context.Context, sessionID string) error
}

type UserRefreshTokenStore interface {
	RefreshTokenStore
	BindSessionToUser(ctx context.Context, userID, sessionID string) error
	RevokeAllByUser(ctx context.Context, userID string) error
}

type SocialIdentity struct {
	Provider       string
	ProviderUserID string
	Email          string
	EmailVerified  bool
	Name           string
	AvatarURL      string
}

type SocialAccountStore interface {
	FindUserBySocial(ctx context.Context, provider, providerUserID string) (userID string, found bool, err error)
	FindOrCreateUserByEmail(ctx context.Context, email string, displayName string) (userID string, err error)
	LinkSocialAccount(ctx context.Context, userID string, identity SocialIdentity) error
}
