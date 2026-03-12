# Cache Module

Genel cache arayüzü (permission, session, tenant, API response için).

## Providers

- `cache/memory` — in-memory (dev/test)
- `cache/redis` — Redis (production)

## Usage

```go
import (
    "github.com/parevo/core/cache"
    "github.com/parevo/core/cache/memory"
)

c := memory.NewCache(5 * time.Minute)
c.Set(ctx, "key", []byte("value"), 10*time.Second)
v, err := c.Get(ctx, "key")
c.Delete(ctx, "key")
```

## Redis

```go
import "github.com/parevo/core/cache/redis"

c := redis.NewCache(redisClient, "app:")
c.Set(ctx, "key", []byte("value"), 10*time.Second)
```
