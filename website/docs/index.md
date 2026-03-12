# Parevo Core

Framework-agnostic Go library for auth, tenant, and permission management.

```bash
go get github.com/parevo/core
```

Import: `github.com/parevo/core/auth`, `github.com/parevo/core/tenant`, etc.

## Core Modules

- **auth** — JWT service, guards, middleware adapters
- **tenant** — tenant selection, override policy, lifecycle
- **permission** — permission check service
- **storage** — DB adapter interfaces
- **notification** — email, SMS, WebSocket
- **blob** — object storage (S3, R2)
- **cache** — generic cache (memory, Redis)
- **health** — health checks (DB, Redis, blob)
- **lock** — distributed lock (memory, Redis)
- **billing** — tenant usage tracking
- **job** — async job queue
- **search** — full-text search (SQL builder)
- **export** — GDPR data export
- **validation** — request/body validation
- **geo** — IP geolocation

## Supported Frameworks

- net/http, chi, gin, echo, fiber
- GraphQL (auth adapter)

## Quick Start

```go
import (
    "github.com/parevo/core/auth"
    "github.com/parevo/core/auth/adapters/nethttp"
)

svc, _ := auth.NewService(auth.Config{
    Issuer:    "parevo",
    Audience:  "parevo-api",
    SecretKey: []byte("your-secret"),
})

mux := http.NewServeMux()
mux.Handle("/secure", nethttp.AuthMiddleware(svc, adapters.Options{})(yourHandler))
```
