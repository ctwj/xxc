// Package semaphore provides a weighted semaphore implementation.
// This is a minimal stub to satisfy dependencies.
package semaphore

import (
	"context"
	"sync"
)

// NewWeighted creates a new weighted semaphore with the given
// maximum combined weight for concurrent access.
func NewWeighted(n int64) *Weighted {
	return &Weighted{size: n}
}

// Weighted provides a way to bound concurrent access to a resource.
type Weighted struct {
	size    int64
	cur     int64
	mu      sync.Mutex
	waiters []chan struct{}
}

// Acquire acquires the semaphore with a weight of n.
func (w *Weighted) Acquire(ctx context.Context, n int64) error {
	w.mu.Lock()
	if w.cur+n <= w.size {
		w.cur += n
		w.mu.Unlock()
		return nil
	}

	// Wait for availability
	ch := make(chan struct{})
	w.waiters = append(w.waiters, ch)
	w.mu.Unlock()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire acquires the semaphore with a weight of n without blocking.
func (w *Weighted) TryAcquire(n int64) bool {
	w.mu.Lock()
	success := w.cur+n <= w.size
	if success {
		w.cur += n
	}
	w.mu.Unlock()
	return success
}

// Release releases the semaphore with a weight of n.
func (w *Weighted) Release(n int64) {
	w.mu.Lock()
	w.cur -= n
	if len(w.waiters) > 0 {
		close(w.waiters[0])
		w.waiters = w.waiters[1:]
	}
	w.mu.Unlock()
}
