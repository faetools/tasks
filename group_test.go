package tasks_test

import (
	"testing"
)

func TestGroup(t *testing.T) {
	t.Parallel()

	// first := tasks.FirstDo(context.Background(), "first", quick)

	// first.ThenDo("another", sleep)
	// first.ThenDo("yet another", sleep)

	// // should take less than double sleep time due to running in parallel
	// assert.Eventually(t, func() bool {
	// 	if err := first.Execute(); err != nil {
	// 		return false
	// 	}
	// 	return true
	// }, sleepTime*3/2, sleepTime/4)
}
