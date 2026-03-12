# MongoDB storage adapter

Production-ready MongoDB adapter for sessions, refresh tokens, tenants, and permissions.

## Setup

MongoDB creates collections automatically on first write. No migration needed. Collections used:

- `parevo_sessions` — session_id, user_id, revoked
- `parevo_refresh_tokens` — token_id, session_id, user_id, replaced_by, expires_at
- `parevo_tenants` — tenant_id, name, owner_id, status
- `parevo_subject_tenants` — subject_id, tenant_id
- `parevo_permission_grants` — subject_id, tenant_id, permission
- `parevo_users` — user_id, email, display_name
- `parevo_social_accounts` — provider, provider_user_id, user_id, email, name, avatar_url

## Usage

```go
import (
    "context"
    "github.com/parevo/core/storage/mongodb"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

ctx := context.Background()
client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
db := client.Database("parevo")

sessionStore := mongodb.NewSessionStore(db)
refreshStore := mongodb.NewRefreshStore(db)
tenantStore := mongodb.NewTenantStore(db)
tenantLifecycleStore := mongodb.NewTenantLifecycleStore(db)
permissionStore := mongodb.NewPermissionStore(db)
socialStore := mongodb.NewSocialAccountStore(db)
```

## Dependencies

```bash
go get go.mongodb.org/mongo-driver/v2
```
