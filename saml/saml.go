package saml

import (
	"context"
	"encoding/xml"
	"errors"
)

var (
	ErrSAMLNotConfigured = errors.New("SAML not configured")
	ErrInvalidSAMLResponse = errors.New("invalid SAML response")
)

// IdentityProvider represents an IdP (e.g. Okta, Azure AD).
type IdentityProvider struct {
	EntityID     string
	SSOURL       string
	Certificate  string
	MetadataURL  string
}

// ServiceProviderConfig holds SP configuration.
type ServiceProviderConfig struct {
	EntityID    string
	ACSURL      string
	MetadataURL string
}

// Service provides SAML 2.0 SSO flows.
type Service struct {
	idp IdentityProvider
	sp  ServiceProviderConfig
}

// NewService creates a SAML service.
func NewService(idp IdentityProvider, sp ServiceProviderConfig) *Service {
	return &Service{idp: idp, sp: sp}
}

// AuthRequestURL returns the URL to redirect the user to for IdP login.
func (s *Service) AuthRequestURL(ctx context.Context, relayState string) (string, error) {
	_ = ctx
	if s.idp.SSOURL == "" {
		return "", ErrSAMLNotConfigured
	}
	// In real impl: build SAML AuthnRequest, encode, return URL
	return s.idp.SSOURL + "?RelayState=" + relayState, nil
}

// ParseResponse parses a SAML response and returns the subject (user identifier).
func (s *Service) ParseResponse(ctx context.Context, rawResponse string) (*SAMLAssertion, error) {
	_ = ctx
	var doc struct {
		XMLName xml.Name `xml:"Response"`
	}
	if err := xml.Unmarshal([]byte(rawResponse), &doc); err != nil {
		return nil, ErrInvalidSAMLResponse
	}
	// In real impl: verify signature, parse NameID/attributes
	return &SAMLAssertion{}, nil
}

// SAMLAssertion holds parsed SAML assertion data.
type SAMLAssertion struct {
	NameID             string
	SessionIndex        string
	Attributes         map[string][]string
}
