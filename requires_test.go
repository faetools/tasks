package tasks_test

import (
	"testing"

	"github.com/faetools/tasks"
	"github.com/stretchr/testify/assert"
)

func TestRequires(t *testing.T) {
	t.Parallel()

	firstStarted := false

	first := tasks.FirstDo("first", func() error {
		firstStarted = true
		return sleep()
	})
	second := tasks.FirstDo("another", func() error {
		if !firstStarted {
			t.Fatal("first didn't start")
		}
		return nil
	}, tasks.OptionRequires(first))

	assert.Never(t, func() bool {
		_ = second.Run()
		return true
	}, sleepTime, sleepTime/4)
}
