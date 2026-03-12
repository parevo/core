package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/parevo/core/blob"
	"github.com/parevo/core/blob/memory"
)

func main() {
	// Option 1: In-memory (dev/test)
	var store blob.Store = memory.NewStore()

	// Option 2: Amazon S3 (has PresignGet, PresignPut)
	// store, _ = s3.NewStore(s3.Config{
	// 	Region:          "us-east-1",
	// 	AccessKeyID:     "your-access-key",
	// 	SecretAccessKey: "your-secret-key",
	// })

	// Option 3: Cloudflare R2 (has PresignGet, PresignPut)
	// store, _ = r2.NewStore(r2.Config{...})

	ctx := context.Background()
	bucket := "demo"
	key := "hello.txt"

	// Create
	if err := store.Put(ctx, bucket, key, bytes.NewReader([]byte("Hello, blob!")), "text/plain"); err != nil {
		panic(err)
	}
	fmt.Println("Put ok")

	// Read
	rc, err := store.Get(ctx, bucket, key)
	if err != nil {
		panic(err)
	}
	defer func() { _ = rc.Close() }()
	data, _ := io.ReadAll(rc)
	fmt.Printf("Get: %s\n", data)

	// List
	infos, err := store.List(ctx, bucket, "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("List: %d objects\n", len(infos))
	for _, info := range infos {
		fmt.Printf("  - %s (%d bytes)\n", info.Key, info.Size)
	}

	// Presigned URL (S3/R2 only; memory does not implement PresignedStore)
	if ps, ok := store.(blob.PresignedStore); ok {
		url, err := ps.PresignGet(ctx, bucket, key, 15*time.Minute)
		if err != nil {
			panic(err)
		}
		fmt.Printf("PresignGet: %s...\n", url[:min(60, len(url))])
	}

	// Delete
	if err := store.Delete(ctx, bucket, key); err != nil {
		panic(err)
	}
	fmt.Println("Delete ok")
}
