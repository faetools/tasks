package tasks

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Task is a piece of work that needs to be done.
type Task struct {
	ctx  context.Context
	name string

	// The task function to run
	f Func

	while    Tasks // what to do in parallel
	requires Tasks // what to do beforehand
	then     Tasks // what to do afterwards

	// some options
	timeout time.Duration
	retries int
	logger  io.Writer

	done bool
	once sync.Once
	err  error
}

func newTask(name string, f Func, opts ...Option) *Task {
	t := &Task{
		ctx:    context.Background(),
		name:   name,
		f:      f,
		logger: io.Discard,
	}

	for _, opt := range opts {
		t.Apply(opt)
	}

	return t
}

// Apply applies some options to the task.
func (t *Task) Apply(o Option) { o(t) }

// Requires specifies a task to be run before this task.
func (t *Task) Requires(required *Task) {
	t.requires = append(t.requires, required)
}

// SetContext sets a context for the task.
func (t *Task) SetContext(ctx context.Context) {
	t.ctx = ctx
}

func (t *Task) runMain() {
	defer func() { t.done = true }()

	f := t.f
	switch {
	case t.f == nil:
		// no function defined, it's a pure grouping of tasks
		f = t.while.Run
	case len(t.while) > 0:
		// run all tasks together
		f = func() error {
			eg := errgroup.Group{}

			eg.Go(t.f)

			for _, w := range t.while {
				w := w
				eg.Go(w.Run)
			}

			return eg.Wait()
		}
	}

	fmt.Fprintf(t.logger, "starting %s...\n", t.name)

	if t.timeout > 0 {
		ctx, cancel := context.WithTimeout(t.ctx, t.timeout)
		defer cancel()
		t.ctx = ctx
	}

	for ; -1 < t.retries; t.retries-- {
		t.err = errors.Wrap(DoWithContext(t.ctx, f), t.name)

		if t.retries > 0 &&
			!errors.Is(t.err, context.DeadlineExceeded) &&
			!errors.Is(t.err, context.Canceled) {
			fmt.Fprintf(t.logger,
				"ERROR trying %s: %s\nmaking another try (%d tries left)...\n",
				t.name, t.err, t.retries-1)
		}
	}

	if t.err != nil {
		fmt.Fprintf(t.logger, "ERROR: %s\n", t.err)
	} else {
		fmt.Fprintf(t.logger, "SUCCESS running %s\n", t.name)
	}
}

// Run executes the task.
// Successive calls will not execute the task again but return the same error.
func (t *Task) Run() error {
	// if the task was already done, return right away
	if t.done {
		return t.err
	}

	// do all required tasks first
	if err := t.requires.Run(); err != nil {
		return err
	}

	// do the tasks itself
	t.once.Do(t.runMain)

	// return error if there was one
	if t.err != nil {
		return t.err
	}

	// run all subtasks
	return t.then.Run()
}
