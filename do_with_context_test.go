package tasks_test

import (
	"context"
	"testing"
	"time"

	"github.com/faetools/tasks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const sleepTime = time.Millisecond * 100

var errDone = errors.New("done")

func quick() error { return nil }

func sleep() error {
	time.Sleep(sleepTime)
	return nil
}

func quickErr() error { return errDone } // nolint:wrapcheck

func sleepErr() error {
	time.Sleep(sleepTime)
	return errDone // nolint:wrapcheck
}

func TestDoWithContext(t *testing.T) {
	t.Parallel()
	bctx := context.Background()

	t.Run("background", func(t *testing.T) {
		t.Parallel()

		err := tasks.DoWithContext(bctx, quickErr)
		assert.ErrorIs(t, err, errDone)
	})

	t.Run("cancel right away", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(bctx)
		cancel()

		assert.ErrorIs(t, tasks.DoWithContext(ctx, quickErr), context.Canceled)
	})

	t.Run("cancel later", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(bctx)

		time.AfterFunc(sleepTime/2, cancel)

		assert.ErrorIs(t, tasks.DoWithContext(ctx, sleep), context.Canceled)
	})

	t.Run("cancel later with timeout", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithDeadline(bctx, time.Now().Add(time.Hour))

		time.AfterFunc(sleepTime/2, cancel)

		assert.ErrorIs(t, tasks.DoWithContext(ctx, sleep), context.Canceled)
	})

	t.Run("with timeout not reached", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithDeadline(bctx, time.Now().Add(time.Hour))

		time.AfterFunc(sleepTime*2, cancel)

		assert.ErrorIs(t, tasks.DoWithContext(ctx, sleepErr), errDone)
	})

	t.Run("with timeout reached", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithDeadline(bctx, time.Now().Add(sleepTime/2))

		time.AfterFunc(sleepTime*2, cancel)

		assert.ErrorIs(t, tasks.DoWithContext(ctx, sleep), context.DeadlineExceeded)
	})
}
