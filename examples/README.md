# Examples

This folder contains working integration examples.

## nethttp-basic

```bash
go run ./examples/nethttp-basic
```

Note: No token endpoint; minimal example showing middleware flow.

## gin-modular

```bash
go run ./examples/gin-modular
```

Shows how to wire `tenant` and `permission` modules with `storage/memory` adapters.

## custom-db-adapter

```bash
go run ./examples/custom-db-adapter
```

Example of mapping a custom DB schema to storage interfaces.

## social-login

```bash
go run ./examples/social-login
```

Shows social callback flow: provider exchange, account linking, access token issuance.

## social-google

```bash
GOOGLE_CLIENT_ID=... GOOGLE_CLIENT_SECRET=... GOOGLE_REDIRECT_URL=... GOOGLE_AUTH_CODE=... go run ./examples/social-google
```

Real Google OIDC provider callback flow.

## refresh-rotation

```bash
go run ./examples/refresh-rotation
```

Refresh token rotation and reuse security flows.

## logging

```bash
go run ./examples/logging
```

For production JSON format:

```bash
PAREVO_ENV=production go run ./examples/logging
```

## apikey

```bash
go run ./examples/apikey
```

API key authentication. Tokens with `pk_` or `sk_` prefix are validated as API keys.

## tenant-ratelimit

```bash
go run ./examples/tenant-ratelimit
```

Tenant-based rate limiting. Use `X-Tenant-Id` header.

## permission-cache

```bash
go run ./examples/permission-cache
```

Cached permission store with TTL and invalidation on role changes.

## totp-mfa

```bash
go run ./examples/totp-mfa
```

TOTP 2FA setup and verification with pquerna/otp.

## admin-panel

```bash
go run ./examples/admin-panel
```

Admin panel at http://localhost:8086/admin. Login required: visit /login, use `admin` / `admin123` (demo). Only users with `admin:*` permission can access. Manage tenants, permissions, sessions.

## notification

```bash
go run ./examples/notification
```

Unified notification: email, SMS, WebSocket via `notification.Sender` interface. Uses `notification/memory` for dev/test.

## blob

```bash
go run ./examples/blob
```

Object storage: Put, Get, List, Delete. Uses `blob/memory` for dev/test. Swap to `blob/s3` or `blob/r2` for production.

## mysql-storage

```bash
# Apply schema first: mysql -u root -p parevo < storage/mysql/schema.sql
MYSQL_DSN="user:pass@tcp(localhost:3306)/parevo?parseTime=true" go run ./examples/mysql-storage
```

MySQL storage adapter: TenantStore, PermissionStore, SessionStore, RefreshStore, SocialAccountStore. Runs on :8083.

## mongodb-storage

```bash
MONGODB_URI="mongodb://localhost:27017" go run ./examples/mongodb-storage
```

MongoDB storage adapter: same stores as MySQL. Collections auto-created. Runs on :8084.
