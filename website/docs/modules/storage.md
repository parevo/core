# Storage Module

DB adapter interfaces for auth, tenant, permission.

## Interfaces

- `TenantStore` — ResolveSubjectTenants
- `PermissionStore` — HasPermission
- `SessionStore` — RevokeSession, IsSessionRevoked
- `RefreshTokenStore` — MarkIssued, IsUsed, MarkUsed

## Implementations

- `storage/memory` — in-memory for dev/test (all stores)
- `storage/postgres` — Postgres (SessionStore, RefreshStore, TenantStore, TenantLifecycleStore, PermissionStore, SocialAccountStore)
- `storage/mysql` — MySQL (same as Postgres)
- `storage/mongodb` — MongoDB (same as Postgres)
- `storage/redis` — Redis (SessionStore, RefreshStore)
