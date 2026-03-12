# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Tenant lifecycle API (`tenant.LifecycleService`): Create, Suspend, Resume, Delete, Status
- `storage.TenantLifecycleStore` interface and `storage/memory` implementation
- Permission wildcard support (`orders:*`, `*:read`, `*:*`) in `storage/memory.PermissionStore`
- `permission.MatchPermission` and `permission.CheckWithWildcard`
- Redis storage adapter (`storage/redis`) for SessionStore and RefreshStore
- GitHub OIDC provider (`social/providers/github`)
- 2FA/MFA (`auth/mfa`): TOTPService, TOTPStore interface
- API key auth (`auth/apikey`): Service, APIKeyStore, GenerateKey
- API key integration: `AuthenticateContext` validates `pk_`/`sk_` prefixed tokens as API keys
- Config strict validation (`auth.Config.ValidateStrict()`)
- Tenant-aware rate limit (`auth/ratelimit.TenantLimiter`): per-tenant QPS, override via `SetTenantLimit`
- `TenantRateLimitMiddleware` (nethttp, chi adapters)
- Permission cache (`permission.CachedPermissionStore`): TTL cache, `InvalidateSubject`, `InvalidateSubjectTenant`, `InvalidateAll`
- TOTP implementation (`auth/mfa.PquernaVerifier`): pquerna/otp for verification and secret generation
- golangci-lint in CI pipeline
- Prometheus metrics, OpenTelemetry tracing, Postgres adapter (previous release)

## [0.1.0] - 2025-03-12

### Added

- JWT auth service with access/refresh token lifecycle
- Key rotation (KID) support
- Refresh token rotation and reuse detection
- Session revoke and logout-all-by-user
- Tenant resolution with override policy
- Permission check service
- Social login callback + Google OIDC provider
- Framework adapters: net/http, chi, gin, echo, fiber
- Rate limiting and brute-force lockout
- Structured logging (dev/prod) with redaction
- Audit logging hooks
- Request ID middleware
- Tenant-safe SQL helper (`auth/tenantsql`)
- In-memory storage adapters
- Examples: nethttp-basic, gin-modular, social-login, social-google, refresh-rotation, logging, custom-db-adapter
