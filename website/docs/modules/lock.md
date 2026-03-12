# Lock Module

Dağıtık lock arayüzü (rate limiting, job deduplication, critical section için).

## Providers

- `lock/memory` — in-memory (tek instance)
- `lock/redis` — Redis (çoklu instance)

## Usage

```go
import (
    "github.com/parevo/core/lock"
    "github.com/parevo/core/lock/memory"
)

l := memory.NewLocker()
ok, err := l.Lock(ctx, "job:123", 30*time.Second)
if !ok {
    return // zaten kilitli
}
defer l.Unlock(ctx, "job:123")
// ... kritik işlem
```

## Redis

```go
import "github.com/parevo/core/lock/redis"

l := redis.NewLocker(redisClient, "lock:")
ok, err := l.Lock(ctx, "job:123", 30*time.Second)
```
