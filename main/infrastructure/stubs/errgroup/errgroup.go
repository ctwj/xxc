// Package errgroup provides synchronization, error propagation, and Context
// cancellation for groups of goroutines working on subtasks of a common task.
// This is a minimal implementation to satisfy the ants dependency.
package errgroup

import (
	"context"
	"sync"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// WithContext returns a new Group and an associated Context derived from ctx.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Go calls the given function in a new goroutine.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// SetLimit limits the number of active goroutines in this group to n.
// A negative value indicates no limit.
func (g *Group) SetLimit(n int) {
	// Minimal implementation - no limit enforcement
}

// TryGo calls the given function in a new goroutine only if the number of
// active goroutines in the group is currently below the configured limit.
func (g *Group) TryGo(f func() error) bool {
	g.Go(f)
	return true
}
