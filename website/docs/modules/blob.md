# Blob Module

Object storage (S3, R2) interface.

## Providers

- `blob/s3` — Amazon S3
- `blob/r2` — Cloudflare R2
- `blob/memory` — dev/test

## Usage

```go
store, _ := s3.NewStore(s3.Config{
    Region: "us-east-1",
    AccessKeyID: "...",
    SecretAccessKey: "...",
})
store.Put(ctx, "bucket", "key", body, "text/plain")
rc, _ := store.Get(ctx, "bucket", "key")
```

## Presigned URLs

S3 and R2 support `PresignGet` and `PresignPut` for direct client access. Memory store implements the interface with placeholder URLs (for interface compatibility in tests).
