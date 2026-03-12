# Examples

Run examples from the project root:

```bash
go run ./examples/nethttp-basic
go run ./examples/gin-modular
go run ./examples/notification
go run ./examples/blob
go run ./examples/admin-panel
# MySQL (requires MySQL + schema): MYSQL_DSN="..." go run ./examples/mysql-storage
# MongoDB: MONGODB_URI="mongodb://localhost:27017" go run ./examples/mongodb-storage
```

## New Modules (no dedicated examples yet)

- **cache** — `cache/memory`, `cache/redis`
- **health** — `health.NewChecker()`, `PingDB`, `PingRedis`, `PingBlob`
- **lock** — `lock/memory`, `lock/redis`
- **billing** — `billing/memory`
- **job** — `job/memory`
- **search** — `search/sql`
- **export** — `export.NewPayload`, `ToJSON`, `ToBlob`
- **validation** — `validation.Validate`, `ValidateJSON`
- **geo** — `geo/memory`

See `examples/README.md` for full list and run instructions.
