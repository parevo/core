package blob

import (
	"context"
	"io"
	"time"
)

// Store provides object storage (S3, R2, etc.).
// CRUD: Put (create/update), Get (read), Delete, List.
type Store interface {
	Put(ctx context.Context, bucket, key string, body io.Reader, contentType string) error
	Get(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, key string) error
	List(ctx context.Context, bucket, prefix string) ([]ObjectInfo, error)
}

// PresignedStore extends Store with presigned URL support for direct client access.
// S3 and R2 implement this; memory does not.
type PresignedStore interface {
	Store
	PresignGet(ctx context.Context, bucket, key string, exp time.Duration) (string, error)
	PresignPut(ctx context.Context, bucket, key string, contentType string, exp time.Duration) (string, error)
}
