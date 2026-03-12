package memory

import (
	"context"
	"strings"
	"time"

	"github.com/parevo/core/storage"
)

type TenantStore struct {
	SubjectTenants map[string][]string
}

func (s *TenantStore) ResolveSubjectTenants(_ context.Context, subjectID string) ([]string, error) {
	if s == nil || s.SubjectTenants == nil {
		return nil, nil
	}
	return s.SubjectTenants[subjectID], nil
}

type TenantLifecycleStore struct {
	Tenants map[string]storage.TenantStatus
}

func (s *TenantLifecycleStore) Create(_ context.Context, tenantID, _, ownerID string) error {
	if s.Tenants == nil {
		s.Tenants = map[string]storage.TenantStatus{}
	}
	if s.Tenants[tenantID] != "" && s.Tenants[tenantID] != storage.TenantStatusDeleted {
		return nil
	}
	s.Tenants[tenantID] = storage.TenantStatusActive
	_ = ownerID
	return nil
}

func (s *TenantLifecycleStore) Suspend(_ context.Context, tenantID string) error {
	if s.Tenants == nil {
		return nil
	}
	s.Tenants[tenantID] = storage.TenantStatusSuspended
	return nil
}

func (s *TenantLifecycleStore) Resume(_ context.Context, tenantID string) error {
	if s.Tenants == nil {
		return nil
	}
	s.Tenants[tenantID] = storage.TenantStatusActive
	return nil
}

func (s *TenantLifecycleStore) Delete(_ context.Context, tenantID string) error {
	if s.Tenants == nil {
		return nil
	}
	s.Tenants[tenantID] = storage.TenantStatusDeleted
	return nil
}

func (s *TenantLifecycleStore) Status(_ context.Context, tenantID string) (storage.TenantStatus, error) {
	if s == nil || s.Tenants == nil {
		return "", nil
	}
	return s.Tenants[tenantID], nil
}

func (s *TenantLifecycleStore) ListTenants(_ context.Context) ([]storage.TenantInfo, error) {
	if s == nil || s.Tenants == nil {
		return nil, nil
	}
	var out []storage.TenantInfo
	for id, status := range s.Tenants {
		if status != storage.TenantStatusDeleted {
			out = append(out, storage.TenantInfo{ID: id, Status: status})
		}
	}
	return out, nil
}

type PermissionStore struct {
	// Key format: subjectID|tenantID|permission
	Grants map[string]bool
}

func (s *PermissionStore) HasPermission(_ context.Context, subjectID, tenantID, permission string, _ []string) (bool, error) {
	if s == nil || s.Grants == nil {
		return false, nil
	}
	prefix := subjectID + "|" + tenantID + "|"
	if s.Grants[prefix+permission] {
		return true, nil
	}
	for k, v := range s.Grants {
		if !v || !strings.HasPrefix(k, prefix) {
			continue
		}
		granted := strings.TrimPrefix(k, prefix)
		if matchPermission(granted, permission) {
			return true, nil
		}
	}
	return false, nil
}

func matchPermission(granted, requested string) bool {
	if granted == requested || granted == "*" || granted == "*:*" {
		return true
	}
	gParts := strings.Split(granted, ":")
	rParts := strings.Split(requested, ":")
	if len(gParts) != 2 || len(rParts) != 2 {
		return false
	}
	if gParts[0] == "*" && gParts[1] == "*" {
		return true
	}
	if gParts[0] == rParts[0] && gParts[1] == "*" {
		return true
	}
	if gParts[0] == "*" && gParts[1] == rParts[1] {
		return true
	}
	return false
}

func (s *PermissionStore) Grant(_ context.Context, subjectID, tenantID, permission string) error {
	if s.Grants == nil {
		s.Grants = map[string]bool{}
	}
	s.Grants[subjectID+"|"+tenantID+"|"+permission] = true
	return nil
}

func (s *PermissionStore) Revoke(_ context.Context, subjectID, tenantID, permission string) error {
	if s.Grants == nil {
		return nil
	}
	delete(s.Grants, subjectID+"|"+tenantID+"|"+permission)
	return nil
}

func (s *PermissionStore) ListGrants(_ context.Context, subjectID, tenantID string) ([]string, error) {
	if s == nil || s.Grants == nil {
		return nil, nil
	}
	prefix := subjectID + "|" + tenantID + "|"
	var out []string
	for k, v := range s.Grants {
		if v && len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, k[len(prefix):])
		}
	}
	return out, nil
}

type SessionStore struct {
	Revoked       map[string]bool
	SessionToUser map[string]string
	UserSessions  map[string]map[string]struct{}
}

func (s *SessionStore) RevokeSession(_ context.Context, sessionID string) error {
	if s.Revoked == nil {
		s.Revoked = map[string]bool{}
	}
	s.Revoked[sessionID] = true
	return nil
}

func (s *SessionStore) IsSessionRevoked(_ context.Context, sessionID string) (bool, error) {
	if s == nil || s.Revoked == nil {
		return false, nil
	}
	return s.Revoked[sessionID], nil
}

func (s *SessionStore) BindSessionToUser(_ context.Context, userID, sessionID string) error {
	if s.SessionToUser == nil {
		s.SessionToUser = map[string]string{}
	}
	if s.UserSessions == nil {
		s.UserSessions = map[string]map[string]struct{}{}
	}
	s.SessionToUser[sessionID] = userID
	if s.UserSessions[userID] == nil {
		s.UserSessions[userID] = map[string]struct{}{}
	}
	s.UserSessions[userID][sessionID] = struct{}{}
	return nil
}

func (s *SessionStore) RevokeAllSessionsByUser(_ context.Context, userID string) error {
	if s.UserSessions == nil {
		return nil
	}
	if s.Revoked == nil {
		s.Revoked = map[string]bool{}
	}
	for sessionID := range s.UserSessions[userID] {
		s.Revoked[sessionID] = true
	}
	return nil
}

func (s *SessionStore) ListSessionsByUser(_ context.Context, userID string) ([]string, error) {
	if s == nil || s.UserSessions == nil {
		return nil, nil
	}
	var out []string
	for sid := range s.UserSessions[userID] {
		out = append(out, sid)
	}
	return out, nil
}

type RefreshStore struct {
	IssuedBySession map[string]map[string]time.Time // session -> tokenID -> expiresAt
	Used            map[string]string               // tokenID -> replacedBy
	UserSessions    map[string]map[string]struct{}  // user -> sessions
}

func (s *RefreshStore) MarkIssued(_ context.Context, sessionID, tokenID string, expiresAt time.Time) error {
	if s.IssuedBySession == nil {
		s.IssuedBySession = map[string]map[string]time.Time{}
	}
	if s.IssuedBySession[sessionID] == nil {
		s.IssuedBySession[sessionID] = map[string]time.Time{}
	}
	s.IssuedBySession[sessionID][tokenID] = expiresAt
	return nil
}

func (s *RefreshStore) IsUsed(_ context.Context, tokenID string) (bool, error) {
	if s == nil || s.Used == nil {
		return false, nil
	}
	_, ok := s.Used[tokenID]
	return ok, nil
}

func (s *RefreshStore) MarkUsed(_ context.Context, tokenID, replacedBy string) error {
	if s.Used == nil {
		s.Used = map[string]string{}
	}
	s.Used[tokenID] = replacedBy
	return nil
}

func (s *RefreshStore) RevokeSession(_ context.Context, sessionID string) error {
	if s.IssuedBySession == nil {
		return nil
	}
	tokens := s.IssuedBySession[sessionID]
	if len(tokens) == 0 {
		return nil
	}
	if s.Used == nil {
		s.Used = map[string]string{}
	}
	for tokenID := range tokens {
		s.Used[tokenID] = "revoked"
	}
	delete(s.IssuedBySession, sessionID)
	return nil
}

func (s *RefreshStore) BindSessionToUser(_ context.Context, userID, sessionID string) error {
	if s.UserSessions == nil {
		s.UserSessions = map[string]map[string]struct{}{}
	}
	if s.UserSessions[userID] == nil {
		s.UserSessions[userID] = map[string]struct{}{}
	}
	s.UserSessions[userID][sessionID] = struct{}{}
	return nil
}

func (s *RefreshStore) RevokeAllByUser(ctx context.Context, userID string) error {
	if s.UserSessions == nil {
		return nil
	}
	for sessionID := range s.UserSessions[userID] {
		if err := s.RevokeSession(ctx, sessionID); err != nil {
			return err
		}
	}
	return nil
}

type UserStore struct {
	Users map[string]storage.UserInfo // userID -> UserInfo
}

func (s *UserStore) ListUsers(_ context.Context) ([]storage.UserInfo, error) {
	if s == nil || s.Users == nil {
		return nil, nil
	}
	out := make([]storage.UserInfo, 0, len(s.Users))
	for _, u := range s.Users {
		out = append(out, u)
	}
	return out, nil
}

type SocialAccountStore struct {
	UsersByProvider map[string]string // provider|providerUserID -> userID
	UsersByEmail    map[string]string // email -> userID
}

func (s *SocialAccountStore) FindUserBySocial(_ context.Context, provider, providerUserID string) (string, bool, error) {
	if s == nil || s.UsersByProvider == nil {
		return "", false, nil
	}
	key := provider + "|" + providerUserID
	userID, ok := s.UsersByProvider[key]
	return userID, ok, nil
}

func (s *SocialAccountStore) FindOrCreateUserByEmail(_ context.Context, email string, _ string) (string, error) {
	if s.UsersByEmail == nil {
		s.UsersByEmail = map[string]string{}
	}
	if userID, ok := s.UsersByEmail[email]; ok {
		return userID, nil
	}
	userID := "user:" + email
	s.UsersByEmail[email] = userID
	return userID, nil
}

func (s *SocialAccountStore) LinkSocialAccount(_ context.Context, userID string, identity storage.SocialIdentity) error {
	if s.UsersByProvider == nil {
		s.UsersByProvider = map[string]string{}
	}
	key := identity.Provider + "|" + identity.ProviderUserID
	s.UsersByProvider[key] = userID
	return nil
}
