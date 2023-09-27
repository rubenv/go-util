package wait

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestWait(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	ctx := context.Background()

	out := make(chan string, 64)

	w := New()

	var g errgroup.Group
	g.Go(func() error {
		assert.NoError(w.Wait(ctx))
		out <- "done"
		return nil
	})
	g.Go(func() error {
		<-w.C
		out <- "done"
		return nil
	})

	assert.Len(out, 0)
	g.Go(func() error {
		out <- "starting"
		w.Broadcast()
		return nil
	})

	err := g.Wait()
	assert.NoError(err)

	assert.Len(out, 3)
	assert.Equal("starting", <-out)
	assert.Equal("done", <-out)
	assert.Equal("done", <-out)

	assert.NoError(w.Wait(ctx)) // Shouldn't hang
	<-w.C                       // Also shouldn't hang
	w.Broadcast()               // No-op
	assert.NoError(w.Wait(ctx)) // Still shouldn't hang
}
