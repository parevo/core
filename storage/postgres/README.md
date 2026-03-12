# Postgres storage adapter

Production-ready Postgres adapter for `SessionStore` and `RefreshStore`.

## Setup

1. Run the migration:

```bash
psql $DATABASE_URL -f schema.sql
```

2. Use in your app:

```go
db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
sessionStore := postgres.NewSessionStore(db)
refreshStore := postgres.NewRefreshStore(db)

svc, _ := auth.NewServiceWithModules(auth.Config{...}, auth.Modules{
    SessionStore: sessionStore,
    RefreshStore: refreshStore,
})
```

## Schema

- `parevo_sessions`: session_id, user_id, revoked
- `parevo_refresh_tokens`: token_id, session_id, user_id, replaced_by, expires_at
