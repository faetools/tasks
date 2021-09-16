package tasks

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var _ UnderlyingTask = (*group)(nil)

// Group are a collection of task groups.
type group struct {
	meta
	collection tasks
}

func newGroup(container *Task, collection tasks) *group {
	return &group{
		meta: meta{
			ctx:       context.Background(),
			logger:    io.Discard,
			container: container,
		},
		collection: collection,
	}
}

func (g *group) WhileDo(name string, f Func, opts ...Option) UnderlyingTask {
	g.collection = append(g.collection, newTask(name, f, opts...))
	return g
}

func (t *group) Run() error {
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
			t.err = errors.Wrap(DoWithContext(t.ctx, t.collection.Run), t.name)

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
