# Postgres storage adapter

Production-ready Postgres adapter for sessions, refresh tokens, tenants, and permissions.

## Setup

1. Run the migration:

```bash
psql $DATABASE_URL -f schema.sql
```

2. Use in your app:

```go
db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))

// Auth
sessionStore := postgres.NewSessionStore(db)
refreshStore := postgres.NewRefreshStore(db)

// Tenant
tenantStore := postgres.NewTenantStore(db)
tenantLifecycleStore := postgres.NewTenantLifecycleStore(db)

// Permission
permissionStore := postgres.NewPermissionStore(db)

// Social (OAuth login)
socialStore := postgres.NewSocialAccountStore(db)

// Admin (optional)
userStore := postgres.NewUserStore(db)
// SessionStore implements SessionListStore; use sessionStore for ListSessionsByUser

// Wire up
tenantSvc := tenant.NewService(tenantStore)
permSvc := permission.NewService(permissionStore)
```

## Schema

| Table | Purpose |
|-------|---------|
| `parevo_sessions` | session_id, user_id, revoked |
| `parevo_refresh_tokens` | token_id, session_id, user_id, replaced_by, expires_at |
| `parevo_tenants` | tenant_id, name, owner_id, status (TenantLifecycleStore) |
| `parevo_subject_tenants` | subject_id, tenant_id (TenantStore.ResolveSubjectTenants) |
| `parevo_permission_grants` | subject_id, tenant_id, permission |
| `parevo_users` | user_id, email, display_name |
| `parevo_social_accounts` | provider, provider_user_id, user_id, email, name, avatar_url |
