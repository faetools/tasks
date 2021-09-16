package tasks

import (
	"golang.org/x/sync/errgroup"
)

// Tasks are a collection of tasks that can be done in parallel.
type Tasks []*Task

// Run executes all tasks in parallel.
func (ts Tasks) Run() error {
	switch len(ts) {
	case 0:
		return nil
	case 1:
		return ts[0].Run()
	}

	eg := errgroup.Group{}

	for _, t := range ts {
		t := t
		eg.Go(t.Run)
	}

	return eg.Wait()
}
