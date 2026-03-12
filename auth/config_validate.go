package auth

import (
	"fmt"
	"strings"
)

func (c Config) ValidateStrict() error {
	var errs []string
	if strings.TrimSpace(c.Issuer) == "" {
		errs = append(errs, "Issuer is required")
	}
	if strings.TrimSpace(c.Audience) == "" {
		errs = append(errs, "Audience is required")
	}
	if len(c.SecretKey) == 0 && len(c.SigningKeys) == 0 {
		errs = append(errs, "SecretKey or SigningKeys is required")
	}
	if len(c.SigningKeys) > 0 {
		if strings.TrimSpace(c.ActiveKID) == "" {
			errs = append(errs, "ActiveKID is required when using SigningKeys")
		} else if len(c.SigningKeys[c.ActiveKID]) == 0 {
			errs = append(errs, "ActiveKID must reference a valid key in SigningKeys")
		}
	}
	if len(c.SecretKey) > 0 && len(c.SecretKey) < 32 {
		errs = append(errs, "SecretKey should be at least 32 bytes for HS256")
	}
	if c.AccessTokenTTL > 0 && c.AccessTokenTTL < 60 {
		errs = append(errs, "AccessTokenTTL should be at least 1 minute")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidConfig, strings.Join(errs, "; "))
	}
	return nil
}
