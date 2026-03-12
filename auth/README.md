# auth package

The `auth` package provides shared JWT infrastructure. Tenant and permission logic can be wired as separate modules.

## Core API

- `NewService(config)`
- `IssueAccessToken(claims)`
- `IssueRefreshToken(claims)`
- `ParseAndValidate(token)`
- `AuthenticateContext(ctx, bearerToken, requestedTenantID, policy)`
- `NewServiceWithModules(config, modules)`
- `AuthorizePermission(ctx, permission)`
- `IssueTokenPair(ctx, claims)`
- `RotateRefreshToken(ctx, refreshToken)`
- `RevokeSession(ctx, sessionID)`
- `RevokeAllSessionsByUser(ctx, userID)`

## Claims

Standard claim set: `sub`, `tenant_id`, `roles`, `permissions`, `session_id`, `typ`

## net/http usage

```go
svc, _ := auth.NewService(auth.Config{
  Issuer: "parevo",
  Audience: "parevo-api",
  SecretKey: []byte("secret"),
})

protected := nethttpadapter.RequestIDMiddleware()(nethttpadapter.AuthMiddleware(svc, adapters.Options{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
})))
```

## Gin / Echo / Fiber

See framework adapters for `AuthMiddleware`, `RequirePermission`, `RateLimitByIP`, `TenantRateLimitMiddleware`.

## Tenant SQL helper

`tenantsql.AppendTenantFilter(query, tenantID, args...)` adds mandatory `tenant_id` filter to queries.

## Modular usage

```go
tenantSvc := tenant.NewService(myTenantStore)
permSvc := permission.NewService(myPermissionStore)

svc, _ := auth.NewServiceWithModules(auth.Config{
  Issuer: "parevo",
  Audience: "parevo-api",
  SecretKey: []byte("secret"),
}, auth.Modules{
  Tenant: tenantSvc,
  Permission: permSvc,
  APIKey: apikey.NewService(apiKeyStore),
})
```

## Production notes

- `NewService` returns error if `Issuer` or `Audience` is empty.
- Use `SigningKeys` and `ActiveKID` for key rotation.
- Default adapter error handler returns generic messages (`unauthorized`, `forbidden`).
- Permission layer defaults to deny when store is not wired.
- With SessionStore, revoked sessions are rejected automatically.
- With RefreshStore, rotation and reuse detection are enabled.
- `auth.NewStructuredAuditLogger` for high-quality audit output.
