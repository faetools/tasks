package tasks_test

import (
	"testing"

	"github.com/faetools/tasks"
	"github.com/stretchr/testify/assert"
)

func TestTaskGraph(t *testing.T) {
	t.Parallel()

	err :=
		tasks.FirstDo("first", quick).
			// ThenDo("second", quick).
			// ThenDo("long third", sleep).
			// WhileDo("error", quickErr).
			// WhileDo("also third", quick).
			// WhileDo("another third", quick).
			Run()

	assert.ErrorIs(t, err, errDone)
}
