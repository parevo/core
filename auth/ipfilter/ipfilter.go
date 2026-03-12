package ipfilter

import (
	"context"
	"errors"
	"net"
)

var (
	ErrIPBlocked = errors.New("IP address blocked")
	ErrIPNotAllowed = errors.New("IP address not in allowlist")
)

// Store provides IP allowlist/blocklist per tenant or global.
type Store interface {
	IsAllowed(ctx context.Context, tenantID, ip string) (bool, error)
	IsBlocked(ctx context.Context, tenantID, ip string) (bool, error)
}

// Service checks IP against allow/block lists.
type Service struct {
	store     Store
	allowMode bool // true: allowlist (deny by default), false: blocklist (allow by default)
}

// NewService creates an IP filter service.
func NewService(store Store, allowlistMode bool) *Service {
	return &Service{store: store, allowMode: allowlistMode}
}

// Allow returns nil if the IP is allowed.
func (s *Service) Allow(ctx context.Context, tenantID, ip string) error {
	if ip == "" {
		return nil
	}
	if net.ParseIP(ip) == nil {
		return nil
	}
	blocked, err := s.store.IsBlocked(ctx, tenantID, ip)
	if err != nil {
		return err
	}
	if blocked {
		return ErrIPBlocked
	}
	if s.allowMode {
		allowed, err := s.store.IsAllowed(ctx, tenantID, ip)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrIPNotAllowed
		}
	}
	return nil
}
