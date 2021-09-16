package tasks_test

import (
	"testing"

	"github.com/faetools/tasks"
	"github.com/stretchr/testify/assert"
)

func TestRetries(t *testing.T) {
	t.Parallel()

	i := 0
	assert.ErrorIs(t, tasks.FirstDo("first", // nolint:paralleltest
		func() error {
			i++
			return errDone
		}, tasks.OptionRetries(3)).Run(), errDone)

	assert.Equal(t, 4, i)

	assert.NoError(t, tasks.FirstDo("", // nolint:paralleltest
		func() error {
			t.Fatal("this code should not be reachable")
			return nil
		}, tasks.OptionRetries(-1)).Run())

	assert.NoError(t, tasks.FirstDo("", // nolint:paralleltest
		func() error {
			t.Fatal("this code should not be reachable")
			return nil
		}, tasks.OptionRetries(-100)).Run())
}
