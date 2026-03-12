<p align="center">
  <img src="website/static/img/logo.svg" alt="Parevo" width="48" height="48" />
</p>

<h1 align="center">Parevo Core</h1>

<p align="center">
  <strong>Framework-agnostic Go library for auth, tenant, and permission management.</strong>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/parevo/core"><img src="https://pkg.go.dev/badge/github.com/parevo/core.svg" alt="Go Reference" /></a>
  <a href="https://github.com/parevo/core/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" /></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go" alt="Go" /></a>
</p>

---

## Features

- **Auth** — JWT, OAuth2, SAML, LDAP, API keys, WebAuthn, magic link
- **Multi-tenant** — Tenant context, lifecycle, feature flags
- **Permission** — RBAC, ABAC, cached checks
- **Storage-agnostic** — MySQL, Postgres, MongoDB, Redis, memory
- **Framework-agnostic** — net/http, chi, gin, echo, fiber, GraphQL

## Quick Start

```bash
go get github.com/parevo/core
```

```go
package main

import (
    "net/http"
    "github.com/parevo/core/auth"
    "github.com/parevo/core/auth/adapters"
    "github.com/parevo/core/auth/adapters/nethttp"
)

func main() {
    svc, _ := auth.NewService(auth.Config{
        Issuer:    "parevo",
        Audience:  "parevo-api",
        SecretKey: []byte("your-secret"),
    })

    mux := http.NewServeMux()
    mux.Handle("/secure", nethttp.AuthMiddleware(svc, adapters.Options{})(yourHandler))
    http.ListenAndServe(":8080", mux)
}
```

## Modules

### Auth & Identity

| Module | Description |
|--------|-------------|
| `auth` | JWT service, guards, middleware adapters |
| `auth/mfa` | TOTP 2FA, recovery codes |
| `auth/apikey` | API key validation |
| `auth/webauthn` | WebAuthn/Passkeys (`-tags webauthn`) |
| `auth/magiclink` | Magic link / email OTP |
| `auth/blacklist` | JWT blacklist for immediate revoke |
| `auth/ipfilter` | IP allowlist/blocklist |
| `auth/oauth2provider` | OAuth2 authorization server |
| `auth/tenantsql` | Tenant filter helpers for SQL |
| `social` | Social login (Google, GitHub) |
| `consent` | OAuth2 consent management |
| `saml` | SAML 2.0 SSO |
| `ldap` | LDAP/Active Directory auth |

### Tenant & Permission

| Module | Description |
|--------|-------------|
| `tenant` | Tenant selection, override policy, lifecycle |
| `tenant/features` | Feature flags, plan limits |
| `permission` | Permission check service |
| `permission/abac` | ABAC conditions |

### Storage & Data

| Module | Description |
|--------|-------------|
| `storage` | DB adapter interfaces |
| `storage/memory` | In-memory adapters |
| `storage/postgres` | Postgres adapter |
| `storage/mysql` | MySQL adapter |
| `storage/mongodb` | MongoDB adapter |
| `storage/redis` | Redis adapter (sessions, refresh) |
| `blob` | Object storage (S3, R2, memory) |
| `cache` | Generic cache (memory, Redis) |
| `lock` | Distributed lock (memory, Redis) |
| `search` | Full-text search (SQL builder) |

### Infrastructure

| Module | Description |
|--------|-------------|
| `health` | Health checks (DB, Redis, blob) |
| `job` | Async job queue (memory) |
| `billing` | Tenant usage tracking |
| `notification` | Email, SMS, WebSocket |
| `webhooks` | Event webhooks |

### Compliance & Utilities

| Module | Description |
|--------|-------------|
| `export` | GDPR data export |
| `validation` | Request/body validation |
| `geo` | IP geolocation |
| `config` | Config validation |
| `observability` | Logging, metrics, tracing, audit |
| `admin` | Admin panel (tenants, permissions, sessions) |

## Supported Frameworks

| Framework | Auth Adapter |
|-----------|--------------|
| net/http | `auth/adapters/nethttp` |
| chi | `auth/adapters/chi` |
| gin | `auth/adapters/gin` |
| echo | `auth/adapters/echo` |
| fiber | `auth/adapters/fiber` |
| GraphQL | `auth/adapters/graphql` |

## Examples

```bash
go run ./examples/nethttp-basic
go run ./examples/gin-modular
go run ./examples/notification
go run ./examples/blob
go run ./examples/admin-panel
```

| Example | Description |
|---------|-------------|
| `nethttp-basic` | Minimal net/http setup |
| `gin-modular` | Auth + tenant + permission |
| `social-login` | Social callback + account linking |
| `totp-mfa` | TOTP 2FA setup and verify |
| `permission-cache` | Cached permission store |
| `tenant-ratelimit` | Tenant-based rate limiting |
| `blacklist-logout` | JWT blacklist on logout |
| `mysql-storage` | MySQL adapter (requires `MYSQL_DSN`) |
| `mongodb-storage` | MongoDB adapter (requires `MONGODB_URI`) |

See [examples/README.md](examples/README.md) for full list and run instructions.

## Documentation

- **Docs:** [parevo.github.io/core](https://parevo.github.io/core/)
- **Local:** `cd website && npm install && npm run start`

## License

MIT. See [LICENSE](LICENSE).

## Contributing

Issue-first workflow. See [.github/CONTRIBUTING.md](.github/CONTRIBUTING.md).
