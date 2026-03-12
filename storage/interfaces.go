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

type TenantInfo struct {
	ID     string
	Status TenantStatus
}

// TenantListStore lists tenants for admin UI. Optional.
type TenantListStore interface {
	ListTenants(ctx context.Context) ([]TenantInfo, error)
}

// PermissionGrantStore grants/revokes permissions for admin UI. Optional.
type PermissionGrantStore interface {
	Grant(ctx context.Context, subjectID, tenantID, permission string) error
	Revoke(ctx context.Context, subjectID, tenantID, permission string) error
	ListGrants(ctx context.Context, subjectID, tenantID string) ([]string, error)
}

// SessionListStore lists sessions for admin UI. Optional.
type SessionListStore interface {
	ListSessionsByUser(ctx context.Context, userID string) ([]string, error)
}

// UserInfo for admin listing.
type UserInfo struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

// UserListStore lists registered users for admin UI. Optional.
type UserListStore interface {
	ListUsers(ctx context.Context) ([]UserInfo, error)
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

// SessionMetadata holds IP, user-agent, last activity for a session.
type SessionMetadata struct {
	SessionID    string
	UserID       string
	IP           string
	UserAgent    string
	LastActivity time.Time
	CreatedAt    time.Time
}

// SessionMetadataStore lists sessions with metadata for admin UI.
type SessionMetadataStore interface {
	SessionStore
	SetMetadata(ctx context.Context, sessionID, userID, ip, userAgent string) error
	UpdateActivity(ctx context.Context, sessionID string) error
	ListWithMetadata(ctx context.Context, userID string) ([]SessionMetadata, error)
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
