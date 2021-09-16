package tasks

import (
	"golang.org/x/sync/errgroup"
)

// Tasks are a collection of tasks that can be done in parallel.
type tasks []*Task

func (ts tasks) WhileDo(name string, f Func, opts ...Option) tasks {
	return append(ts, newTask(name, f, opts...))
}

// Run executes all tasks in parallel.
func (ts tasks) Run() error {
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
