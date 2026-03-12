package memory

import (
	"context"
	"sync"

	"github.com/parevo/core/job"
)

type jobItem struct {
	queue   string
	payload []byte
}

// Queue implements job.Queue and job.Runner with in-memory channels.
type Queue struct {
	mu     sync.Mutex
	queues map[string]chan []byte
	buf    int
}

// NewQueue creates an in-memory job queue. buf is the channel buffer size per queue.
func NewQueue(buf int) *Queue {
	if buf <= 0 {
		buf = 100
	}
	return &Queue{
		queues: make(map[string]chan []byte),
		buf:    buf,
	}
}

// Enqueue adds a job to the queue.
func (q *Queue) Enqueue(ctx context.Context, queueName string, payload []byte) error {
	q.mu.Lock()
	ch, ok := q.queues[queueName]
	if !ok {
		ch = make(chan []byte, q.buf)
		q.queues[queueName] = ch
	}
	q.mu.Unlock()

	select {
	case ch <- append([]byte(nil), payload...):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Run starts a worker that consumes from the queue. Blocks until ctx is cancelled.
func (q *Queue) Run(ctx context.Context, queueName string, handler job.Handler) error {
	q.mu.Lock()
	ch, ok := q.queues[queueName]
	if !ok {
		ch = make(chan []byte, q.buf)
		q.queues[queueName] = ch
	}
	q.mu.Unlock()

	for {
		select {
		case payload := <-ch:
			_ = handler(ctx, payload)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

var _ job.Queue = (*Queue)(nil)
var _ job.Runner = (*Queue)(nil)
