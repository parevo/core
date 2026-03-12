# Storage Module

DB adapter interfaces for auth, tenant, permission.

## Interfaces

- `TenantStore` — ResolveSubjectTenants
- `PermissionStore` — HasPermission
- `SessionStore` — RevokeSession, IsSessionRevoked
- `RefreshTokenStore` — MarkIssued, IsUsed, MarkUsed

## Implementations

- `storage/memory` — in-memory for dev/test
- `storage/postgres` — Postgres (SessionStore, RefreshStore)
- `storage/redis` — Redis (SessionStore, RefreshStore)
