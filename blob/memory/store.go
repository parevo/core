package memory

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/parevo/core/blob"
)

// Store implements blob.Store with in-memory storage for dev/test.
type Store struct {
	mu      sync.RWMutex
	objects map[string][]byte // bucket|key -> content
	meta    map[string]blob.ObjectInfo
}

// NewStore creates an in-memory blob store.
func NewStore() *Store {
	return &Store{
		objects: make(map[string][]byte),
		meta:    make(map[string]blob.ObjectInfo),
	}
}

func key(bucket, key string) string {
	return bucket + "|" + key
}

// Put uploads an object.
func (s *Store) Put(ctx context.Context, bucket, k string, body io.Reader, contentType string) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("memory put read: %w", err)
	}
	s.mu.Lock()
	s.objects[key(bucket, k)] = data
	s.meta[key(bucket, k)] = blob.ObjectInfo{
		Key:          k,
		Size:         int64(len(data)),
		LastModified: time.Now(),
		ContentType:  contentType,
	}
	s.mu.Unlock()
	return nil
}

// Get downloads an object.
func (s *Store) Get(ctx context.Context, bucket, k string) (io.ReadCloser, error) {
	s.mu.RLock()
	data, ok := s.objects[key(bucket, k)]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("blob: object not found: %s/%s", bucket, k)
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// Delete removes an object.
func (s *Store) Delete(ctx context.Context, bucket, k string) error {
	s.mu.Lock()
	delete(s.objects, key(bucket, k))
	delete(s.meta, key(bucket, k))
	s.mu.Unlock()
	return nil
}

// List returns objects with the given prefix.
func (s *Store) List(ctx context.Context, bucket, prefix string) ([]blob.ObjectInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	bucketPrefix := bucket + "|"
	var out []blob.ObjectInfo
	for k, info := range s.meta {
		if !strings.HasPrefix(k, bucketPrefix) {
			continue
		}
		keyPart := k[len(bucketPrefix):]
		if prefix == "" || strings.HasPrefix(keyPart, prefix) {
			out = append(out, info)
		}
	}
	return out, nil
}

var _ blob.Store = (*Store)(nil)
