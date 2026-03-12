# Parevo Core

Framework-agnostic Go library for auth, tenant, and permission management.

```bash
go get github.com/parevo/core
```

Import: `github.com/parevo/core/auth`, `github.com/parevo/core/tenant`, etc.

Core modules: `auth`, `tenant`, `permission`, and `storage` — JWT-based authentication, tenant context, and permission checks as separate composable services.

## Modules

- `auth`: JWT service, guards, middleware adapters
- `auth/mfa`: 2FA/MFA TOTP (pquerna/otp)
- `auth/apikey`: API key validation
- `auth/tenantsql`: tenant filter helpers
- `tenant`: tenant selection, override policy, lifecycle (create/suspend/delete)
- `permission`: permission check service
- `social`: social login callback and account linking
- `social/providers/google`: Google OIDC
- `social/providers/github`: GitHub OAuth
- `storage`: DB adapter interfaces
- `storage/memory`: in-memory adapters for quick start
- `storage/postgres`: Postgres adapter (SessionStore, RefreshStore)
- `storage/redis`: Redis adapter (SessionStore, RefreshStore)
- `observability/logging`: structured logging (dev/prod)
- `observability/metrics`: Prometheus metrics
- `observability/tracing`: OpenTelemetry-compatible tracing

## Supported Frameworks

- `net/http`
- `chi` (on top of net/http)
- `gin`
- `echo`
- `fiber`

## Goal

Reuse auth + tenant semantics across projects without code duplication.

## Product Readiness

- Config requirements: `Issuer`, `Audience`, `SecretKey`
- Key rotation: `SigningKeys` + `ActiveKID`
- Refresh rotation + reuse detection
- Session revoke + refresh-store integration
- Logout-all (user session family revoke)
- Sanitized error responses (no detail leakage)
- Rate limit middleware (IP and tenant-aware)
- Brute-force lockout manager
- Structured audit logging (dev/prod format + redaction)
- Framework parity tests: net/http, gin, echo, fiber, chi
- Tenant filtering: `auth/tenantsql` for mandatory `tenant_id`
- Modular DB adapter model via `storage` interfaces

## Examples

- `examples/nethttp-basic`: minimal net/http setup
- `examples/gin-modular`: auth + tenant + permission integration
- `examples/custom-db-adapter`: mapping custom DB schema to storage interfaces
- `examples/social-login`: social callback + link + token
- `examples/social-google`: real Google OIDC provider
- `examples/refresh-rotation`: refresh rotation security flow
- `examples/logging`: dev/prod log formats
- `examples/apikey`: API key authentication
- `examples/tenant-ratelimit`: tenant-based rate limiting
- `examples/permission-cache`: cached permission store with invalidation
- `examples/totp-mfa`: TOTP 2FA setup and verify

See `examples/README.md` for run instructions.

## Release

- `CHANGELOG.md`: semantic versioning
- `Makefile`: test, vet, fmt, lint
- `.github/workflows/ci.yml`: CI pipeline

## License

MIT. See [LICENSE](LICENSE).

## Contributing

Issue-first workflow. See `.github/CONTRIBUTING.md` and open issues for feature requests.

## Structure

See `STRUCTURE.md` for folder layout.
