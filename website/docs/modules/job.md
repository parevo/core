# Job Module

Async task/job queue arayüzü.

## Providers

- `job/memory` — in-memory (dev/tek instance)

## Usage

```go
import (
    "github.com/parevo/core/job"
    "github.com/parevo/core/job/memory"
)

queue := memory.NewQueue(100)
queue.Enqueue(ctx, "email", []byte(`{"to":"a@b.com","subject":"Hi"}`))

// Worker
go queue.Run(ctx, "email", func(ctx context.Context, payload []byte) error {
    var msg EmailPayload
    json.Unmarshal(payload, &msg)
    return sendEmail(msg)
})
```
