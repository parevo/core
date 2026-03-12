package auth

import "time"

type Config struct {
	Issuer          string
	Audience        string
	SecretKey       []byte
	ActiveKID       string
	SigningKeys     map[string][]byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Leeway          time.Duration
}

func (c Config) withDefaults() Config {
	cfg := c
	if cfg.AccessTokenTTL <= 0 {
		cfg.AccessTokenTTL = 15 * time.Minute
	}
	if cfg.RefreshTokenTTL <= 0 {
		cfg.RefreshTokenTTL = 7 * 24 * time.Hour
	}
	if cfg.Leeway < 0 {
		cfg.Leeway = 0
	}
	return cfg
}
