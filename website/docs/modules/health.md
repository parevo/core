# Health Module

DB, Redis, blob ve özel kontroller için health check.

## Usage

```go
import (
    "github.com/parevo/core/health"
)

h := health.NewChecker()
h.Add("db", health.PingDB(db))
h.Add("redis", health.PingRedis(redisClient))
h.Add("s3", health.PingBlob(blobStore, "bucket", ""))

// Basit kontrol
if !h.Check(ctx) {
    w.WriteHeader(http.StatusServiceUnavailable)
}

// Detaylı sonuçlar
results := h.CheckWithResults(ctx)
for name, err := range results {
    if err != nil {
        log.Printf("%s: %v", name, err)
    }
}
```

## Blob Kontrolü

- `PingBlob` — List çağrısı ile erişilebilirlik
- `PingBlobPutGet` — Yazma/okuma/ silme ile daha kapsamlı test
