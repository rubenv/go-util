package wait

import (
	"context"
	"sync"
)

type Wait struct {
	o sync.Once
	C chan interface{}
}

func New() *Wait {
	return &Wait{
		C: make(chan interface{}),
	}
}

func (w *Wait) Wait(ctx context.Context) error {
	select {
	case <-w.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *Wait) Broadcast() {
	w.o.Do(func() {
		close(w.C)
	})
}
