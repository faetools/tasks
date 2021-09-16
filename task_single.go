package tasks

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

var _ UnderlyingTask = (*simpleTask)(nil)

// A Task is a single piece of work.
type simpleTask struct {
	// The task function to be run
	f Func

	// meta data
	meta
}

func (t *simpleTask) WhileDo(name string, f Func, opts ...Option) UnderlyingTask {
	return newGroup(t.container, tasks{t.container, newTask(name, f, opts...)})
}

// Run executes the task.
// Successive calls will not execute the task again but return the same error.
func (t *simpleTask) Run() error {
	// if the task was already done, return right away
	if t.done {
		return t.err
	}

	// do all required tasks first
	if err := t.requires.Run(); err != nil {
		return err
	}

	// do the tasks itself
	t.once.Do(func() {
		fmt.Fprintf(t.logger, "starting %s...", t.name)

		if t.timeout > 0 {
			var cancel context.CancelFunc
			t.ctx, cancel = context.WithTimeout(t.ctx, t.timeout)
			defer cancel()
		}

		for ; -1 < t.retries; t.retries-- {
			t.err = errors.Wrap(DoWithContext(t.ctx, t.f), t.name)

			if t.retries > 0 &&
				!errors.Is(t.err, context.DeadlineExceeded) &&
				!errors.Is(t.err, context.Canceled) {
				fmt.Fprintf(t.logger,
					"ERROR trying %s: %s\nmaking another try (%d tries left)...",
					t.name, t.err, t.retries-1)
			}
		}
	})

	// return error if there was one
	if t.err != nil {
		return t.err
	}

	// run all subtasks
	return t.then.Run()
}
