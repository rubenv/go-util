package wait

import (
	"context"
	"sync"
)

// Provides a synchronization point for waiting on a condition to happen (once)
//
// Every call to Wait() will block until Broadcast() is called, then continue afterwards.
type Wait struct {
	o sync.Once
	C chan interface{}
}

func New() *Wait {
	return &Wait{
		C: make(chan interface{}),
	}
}

// Wait for Broadcast() to be called.
//
// Returns nil when Broadcast() was called, a context error if the context was aborted.
func (w *Wait) Wait(ctx context.Context) error {
	select {
	case <-w.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Signal all waiting go-routines that they can continue.
//
// Can be called multiple times.
func (w *Wait) Broadcast() {
	w.o.Do(func() {
		close(w.C)
	})
}
