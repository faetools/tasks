package tasks_test

import (
	"os"
	"testing"

	"github.com/faetools/tasks"
	"github.com/stretchr/testify/assert"
)

func TestParallel(t *testing.T) {
	t.Parallel()

	first := tasks.FirstDo("first", quick, tasks.OptionLogger(os.Stdout))

	first.ThenDo("another", sleep)
	first.ThenDo("yet another", sleep)

	// should take less than double sleep time due to running in parallel
	assert.Eventually(t, func() bool {
		if err := first.Run(); err != nil {
			return false
		}
		return true
	}, sleepTime*5/3, sleepTime/4)
}
