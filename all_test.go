package tasks_test

import (
	"os"
	"testing"

	"github.com/faetools/tasks"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Parallel()

	err :=
		tasks.FirstDo("first", quick, tasks.OptionLogger(os.Stdout)).
			ThenDo("second", quick, tasks.OptionLogger(os.Stdout)).
			ThenDo("long third", sleep, tasks.OptionLogger(os.Stdout)).
			WhileDo("error", quickErr, tasks.OptionLogger(os.Stdout)).
			WhileDo("also third", quick, tasks.OptionLogger(os.Stdout)).
			WhileDo("another third", quick, tasks.OptionLogger(os.Stdout)).
			Run()

	assert.ErrorIs(t, err, errDone)

	err = tasks.FirstDo("container", nil, tasks.OptionLogger(os.Stdout)).
		WhileDo("error", quickErr, tasks.OptionLogger(os.Stdout)).
		Run()
	assert.EqualError(t, err, "container: error: done")
}
