// Package job provides async task/job queue interface for background processing.
//
// # Usage
//
//	queue := memory.NewQueue()
//	queue.Enqueue(ctx, "email", []byte(`{"to":"a@b.com"}`))
//	queue.Run(ctx, "email", func(payload []byte) error { ... })
package job

import (
	"context"
)

// Queue enqueues and processes jobs.
type Queue interface {
	Enqueue(ctx context.Context, queueName string, payload []byte) error
}

// Handler processes a job payload.
type Handler func(ctx context.Context, payload []byte) error

// Runner runs workers that consume from a queue.
type Runner interface {
	Run(ctx context.Context, queueName string, handler Handler) error
}
