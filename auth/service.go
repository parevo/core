package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parevo/core/storage"
	tenantpkg "github.com/parevo/core/tenant"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Service struct {
	cfg       Config
	modules   Modules
	activeKID string
}

func NewService(cfg Config) (*Service, error) {
	if strings.TrimSpace(cfg.Issuer) == "" {
		return nil, fmt.Errorf("%w: Issuer is required", ErrInvalidConfig)
	}
	if strings.TrimSpace(cfg.Audience) == "" {
		return nil, fmt.Errorf("%w: Audience is required", ErrInvalidConfig)
	}
	activeKID := strings.TrimSpace(cfg.ActiveKID)
	if len(cfg.SigningKeys) > 0 {
		if activeKID == "" {
			return nil, fmt.Errorf("%w: ActiveKID is required when SigningKeys are used", ErrInvalidConfig)
		}
		if len(cfg.SigningKeys[activeKID]) == 0 {
			return nil, fmt.Errorf("%w: active key is missing in SigningKeys", ErrInvalidConfig)
		}
	} else if len(cfg.SecretKey) == 0 {
		return nil, fmt.Errorf("%w: SecretKey or SigningKeys is required", ErrInvalidConfig)
	} else if activeKID == "" {
		activeKID = "legacy-hs256"
	}
	cfg = cfg.withDefaults()
	return &Service{cfg: cfg, activeKID: activeKID}, nil
}

func NewServiceWithModules(cfg Config, modules Modules) (*Service, error) {
	s, err := NewService(cfg)
	if err != nil {
		return nil, err
	}
	s.modules = modules
	return s, nil
}

func (s *Service) IssueAccessToken(base Claims) (string, error) {
	return s.issueToken(base, TokenTypeAccess, s.cfg.AccessTokenTTL)
}

func (s *Service) IssueRefreshToken(base Claims) (string, error) {
	return s.issueRefreshToken(base)
}

func (s *Service) ParseAndValidate(tokenString string) (*Claims, error) {
	claims := &Claims{}

	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("%w: unexpected signing method", ErrInvalidToken)
		}
		if len(s.cfg.SigningKeys) > 0 {
			kid, _ := token.Header["kid"].(string)
			key, ok := s.cfg.SigningKeys[kid]
			if !ok || len(key) == 0 {
				return nil, fmt.Errorf("%w: unknown key id", ErrInvalidToken)
			}
			return key, nil
		}
		return s.cfg.SecretKey, nil
	}, jwt.WithLeeway(s.cfg.Leeway), jwt.WithAudience(s.cfg.Audience), jwt.WithIssuer(s.cfg.Issuer))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !parsed.Valid {
		return nil, ErrInvalidToken
	}

	if claims.UserID == "" {
		return nil, fmt.Errorf("%w: missing sub", ErrInvalidToken)
	}
	if claims.SessionID != "" && s.modules.SessionStore != nil {
		revoked, revokeErr := s.modules.SessionStore.IsSessionRevoked(context.Background(), claims.SessionID)
		if revokeErr != nil {
			return nil, revokeErr
		}
		if revoked {
			return nil, ErrSessionRevoked
		}
	}
	if claims.TokenType == TokenTypeRefresh && claims.ID != "" && s.modules.RefreshStore != nil {
		used, useErr := s.modules.RefreshStore.IsUsed(context.Background(), claims.ID)
		if useErr != nil {
			return nil, useErr
		}
		if used {
			return nil, ErrRefreshReuse
		}
	}

	return claims, nil
}

func (s *Service) AuthenticateContext(ctx context.Context, bearerToken string, requestedTenantID string, policy TenantOverridePolicy) (context.Context, *Claims, error) {
	if s.modules.Tracer != nil {
		var end func()
		ctx, end = s.modules.Tracer.StartSpan(ctx, "auth.authenticate", nil)
		defer end()
	}
	token := extractBearerToken(bearerToken)
	if token == "" {
		s.audit(ctx, AuditAuthFailed, map[string]string{"reason": "missing_token"})
		return ctx, nil, ErrUnauthenticated
	}

	if s.modules.APIKey != nil && (strings.HasPrefix(token, "pk_") || strings.HasPrefix(token, "sk_")) {
		userID, tenantID, err := s.modules.APIKey.Validate(ctx, token)
		if err != nil {
			s.audit(ctx, AuditAuthFailed, map[string]string{"reason": "invalid_apikey"})
			return ctx, nil, ErrUnauthenticated
		}
		if userID == "" {
			s.audit(ctx, AuditAuthFailed, map[string]string{"reason": "invalid_apikey"})
			return ctx, nil, ErrUnauthenticated
		}
		resolvedTenant := tenantID
		if requestedTenantID != "" && requestedTenantID != tenantID {
			if policy == nil || !policy.CanOverride(&Claims{UserID: userID}, requestedTenantID) {
				s.audit(ctx, AuditTenantMismatch, map[string]string{"user_id": userID})
				return ctx, nil, ErrTenantMismatch
			}
			resolvedTenant = requestedTenantID
		}
		claims := &Claims{UserID: userID, TenantID: resolvedTenant}
		s.audit(ctx, AuditAuthSuccess, map[string]string{"user_id": userID, "tenant_id": resolvedTenant, "flow": "apikey"})
		return WithClaims(ctx, claims), claims, nil
	}

	claims, err := s.ParseAndValidate(token)
	if err != nil {
		s.audit(ctx, AuditAuthFailed, map[string]string{"reason": "invalid_token"})
		return ctx, nil, ErrUnauthenticated
	}

	tenantID, err := s.resolveTenant(ctx, claims, requestedTenantID, policy)
	if err != nil {
		if errors.Is(err, ErrTenantMismatch) {
			s.audit(ctx, AuditTenantMismatch, map[string]string{"user_id": claims.UserID})
		}
		return ctx, nil, err
	}
	cloned := *claims
	cloned.TenantID = tenantID
	s.audit(ctx, AuditAuthSuccess, map[string]string{"user_id": cloned.UserID, "tenant_id": cloned.TenantID})
	return WithClaims(ctx, &cloned), &cloned, nil
}

func extractBearerToken(headerValue string) string {
	parts := strings.Fields(strings.TrimSpace(headerValue))
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	if strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func (s *Service) resolveTenant(ctx context.Context, claims *Claims, requestedTenantID string, legacyPolicy TenantOverridePolicy) (string, error) {
	tenantID := claims.TenantID
	if s.modules.Tenant == nil {
		if requestedTenantID != "" && requestedTenantID != tenantID {
			if legacyPolicy == nil || !legacyPolicy.CanOverride(claims, requestedTenantID) {
				return "", ErrTenantMismatch
			}
			return requestedTenantID, nil
		}
		return tenantID, nil
	}

	subject := storage.Subject{
		ID:    claims.UserID,
		Roles: claims.Roles,
	}

	var override tenantpkg.OverridePolicy
	if s.modules.TenantOverride != nil {
		override = s.modules.TenantOverride
	}
	resolved, err := s.modules.Tenant.Resolve(ctx, subject, claims.TenantID, requestedTenantID, override)
	if errors.Is(err, tenantpkg.ErrNoTenant) {
		return "", ErrMissingTenant
	}
	if errors.Is(err, tenantpkg.ErrAccessDenied) {
		return "", ErrTenantMismatch
	}
	return resolved, err
}

func HTTPStatusForError(err error) int {
	switch {
	case err == nil:
		return 200
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrTenantMismatch):
		return 403
	case errors.Is(err, ErrUnauthenticated), errors.Is(err, ErrInvalidToken), errors.Is(err, ErrSessionRevoked), errors.Is(err, ErrInvalidRefresh), errors.Is(err, ErrRefreshReuse):
		return 401
	case errors.Is(err, ErrMissingTenant), errors.Is(err, ErrInvalidConfig):
		return 400
	default:
		return 500
	}
}

func (s *Service) RevokeSession(ctx context.Context, sessionID string) error {
	if s.modules.SessionStore == nil {
		return nil
	}
	if err := s.modules.SessionStore.RevokeSession(ctx, sessionID); err != nil {
		return err
	}
	s.audit(ctx, AuditSessionRevoked, map[string]string{"session_id": sessionID})
	return nil
}

func (s *Service) RevokeAllSessionsByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("%w: userID is required", ErrInvalidConfig)
	}
	if sessionStore, ok := s.modules.SessionStore.(storage.UserSessionStore); ok {
		if err := sessionStore.RevokeAllSessionsByUser(ctx, userID); err != nil {
			return err
		}
	}
	if refreshStore, ok := s.modules.RefreshStore.(storage.UserRefreshTokenStore); ok {
		if err := refreshStore.RevokeAllByUser(ctx, userID); err != nil {
			return err
		}
	}
	s.audit(ctx, AuditSessionRevoked, map[string]string{"user_id": userID, "scope": "all"})
	return nil
}

func (s *Service) issueToken(base Claims, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		UserID:      base.UserID,
		TenantID:    base.TenantID,
		Roles:       base.Roles,
		Permissions: base.Permissions,
		SessionID:   base.SessionID,
		TokenType:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        base.ID,
			Issuer:    s.cfg.Issuer,
			Audience:  []string{s.cfg.Audience},
			Subject:   base.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if s.activeKID != "" {
		token.Header["kid"] = s.activeKID
	}
	signingKey := s.cfg.SecretKey
	if len(s.cfg.SigningKeys) > 0 {
		signingKey = s.cfg.SigningKeys[s.activeKID]
	}
	signed, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (s *Service) IssueTokenPair(ctx context.Context, base Claims) (TokenPair, error) {
	if base.UserID == "" {
		return TokenPair{}, ErrInvalidConfig
	}
	if base.SessionID == "" {
		base.SessionID = mustTokenID()
	}
	if sessionStore, ok := s.modules.SessionStore.(storage.UserSessionStore); ok {
		if err := sessionStore.BindSessionToUser(ctx, base.UserID, base.SessionID); err != nil {
			return TokenPair{}, err
		}
	}
	if refreshStore, ok := s.modules.RefreshStore.(storage.UserRefreshTokenStore); ok {
		if err := refreshStore.BindSessionToUser(ctx, base.UserID, base.SessionID); err != nil {
			return TokenPair{}, err
		}
	}
	accessToken, err := s.IssueAccessToken(base)
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, err := s.issueRefreshToken(base)
	if err != nil {
		return TokenPair{}, err
	}
	s.audit(ctx, AuditAuthSuccess, map[string]string{"user_id": base.UserID, "tenant_id": base.TenantID})
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Service) RotateRefreshToken(ctx context.Context, refreshToken string) (TokenPair, error) {
	claims, err := s.ParseAndValidate(refreshToken)
	if err != nil {
		if errors.Is(err, ErrRefreshReuse) {
			s.audit(ctx, AuditAuthFailed, map[string]string{"reason": "refresh_reuse"})
			return TokenPair{}, ErrRefreshReuse
		}
		return TokenPair{}, ErrInvalidRefresh
	}
	if claims.TokenType != TokenTypeRefresh || claims.ID == "" || claims.SessionID == "" {
		return TokenPair{}, ErrInvalidRefresh
	}
	if s.modules.RefreshStore == nil {
		return TokenPair{}, fmt.Errorf("%w: RefreshStore is required", ErrInvalidConfig)
	}

	used, err := s.modules.RefreshStore.IsUsed(ctx, claims.ID)
	if err != nil {
		return TokenPair{}, err
	}
	if used {
		_ = s.RevokeSession(ctx, claims.SessionID)
		_ = s.modules.RefreshStore.RevokeSession(ctx, claims.SessionID)
		s.audit(ctx, AuditAuthFailed, map[string]string{"reason": "refresh_reuse", "session_id": claims.SessionID})
		return TokenPair{}, ErrRefreshReuse
	}

	base := Claims{
		UserID:      claims.UserID,
		TenantID:    claims.TenantID,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		SessionID:   claims.SessionID,
	}

	accessToken, err := s.IssueAccessToken(base)
	if err != nil {
		return TokenPair{}, err
	}
	newRefreshID := mustTokenID()
	base.ID = newRefreshID
	newRefresh, err := s.issueRefreshToken(base)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.modules.RefreshStore.MarkUsed(ctx, claims.ID, newRefreshID); err != nil {
		return TokenPair{}, err
	}
	s.audit(ctx, AuditAuthSuccess, map[string]string{"user_id": claims.UserID, "tenant_id": claims.TenantID, "flow": "refresh_rotate"})
	return TokenPair{AccessToken: accessToken, RefreshToken: newRefresh}, nil
}

func (s *Service) issueRefreshToken(base Claims) (string, error) {
	if base.ID == "" {
		base.ID = mustTokenID()
	}
	token, err := s.issueToken(base, TokenTypeRefresh, s.cfg.RefreshTokenTTL)
	if err != nil {
		return "", err
	}
	if s.modules.RefreshStore != nil && base.SessionID != "" {
		if markErr := s.modules.RefreshStore.MarkIssued(context.Background(), base.SessionID, base.ID, time.Now().UTC().Add(s.cfg.RefreshTokenTTL)); markErr != nil {
			return "", markErr
		}
	}
	return token, nil
}

func mustTokenID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}
