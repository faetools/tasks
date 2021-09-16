package tasks_test

import (
	"context"
	"os"
	"testing"

	"github.com/faetools/tasks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	t.Parallel()
	bctx := context.Background()

	ctx, cancel := context.WithTimeout(bctx, sleepTime/2)
	defer cancel()

	task := tasks.FirstDo("first", sleep,
		tasks.OptionLogger(os.Stdout),
		tasks.OptionContext(ctx))

	assert.Eventually(t, func() bool {
		// nolint:paralleltest
		return errors.Is(task.Run(), context.DeadlineExceeded)
	}, sleepTime*2/3, sleepTime/8)

	task = tasks.FirstDo("first", sleep,
		tasks.OptionLogger(os.Stdout),
		tasks.OptionTimeout(sleepTime/2))

	assert.Eventually(t, func() bool {
		// nolint:paralleltest
		return errors.Is(task.Run(), context.DeadlineExceeded)
	}, sleepTime*2/3, sleepTime/8)
}
