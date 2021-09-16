package tasks

import (
	"context"
	"io"
)

// Task is a piece of work that needs to be done.
type Task struct {
	t UnderlyingTask
}

func newTask(name string, f Func, opts ...Option) *Task {
	simple := &simpleTask{
		meta: meta{
			ctx:    context.Background(),
			name:   name,
			logger: io.Discard,
		},
		f: f,
	}

	t := &Task{t: simple}
	simple.container = t

	for _, opt := range opts {
		t.Apply(opt)
	}

	return t
}

// Apply applies some options to the task.
// The type of the underlying task might be changed.
func (t *Task) Apply(o Option) {
	t.t = o(t.t)
}

// Requires specifies a task to be run before this task.
func (t *Task) Requires(required *Task) {
	t.t = t.t.Requires(required)
}

// Run executes the task.
func (t *Task) Run() error {
	return t.t.Run()
}
