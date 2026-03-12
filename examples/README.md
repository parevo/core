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
