# MySQL storage adapter

Production-ready MySQL adapter for sessions, refresh tokens, tenants, and permissions.

## Setup

1. Run the migration:

```bash
mysql -u user -p database < schema.sql
```

2. Use in your app:

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    mysqlstorage "github.com/parevo/core/storage/mysql"
)

db, _ := sql.Open("mysql", "user:pass@tcp(localhost:3306)/parevo?parseTime=true")

sessionStore := mysqlstorage.NewSessionStore(db)
refreshStore := mysqlstorage.NewRefreshStore(db)
tenantStore := mysqlstorage.NewTenantStore(db)
tenantLifecycleStore := mysqlstorage.NewTenantLifecycleStore(db)
permissionStore := mysqlstorage.NewPermissionStore(db)
socialStore := mysqlstorage.NewSocialAccountStore(db)
```

## Dependencies

```bash
go get github.com/go-sql-driver/mysql
```

## Schema

Same structure as Postgres: `parevo_sessions`, `parevo_refresh_tokens`, `parevo_tenants`, `parevo_subject_tenants`, `parevo_permission_grants`, `parevo_users`, `parevo_social_accounts`.
