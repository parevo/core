package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID      string   `json:"sub"`
	TenantID    string   `json:"tenant_id"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scopes      []string `json:"scope,omitempty"` // OAuth2 scopes, e.g. "read:orders write:users"
	SessionID   string   `json:"session_id,omitempty"`
	TokenType   string   `json:"typ,omitempty"`
	jwt.RegisteredClaims
}
