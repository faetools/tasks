package tasks

import "context"

// UnderlyingTask is a piece of work that needs to be done.
type UnderlyingTask interface {
	// Apply applies some options to the task.
	// Returns the updated task.
	Apply(Option) UnderlyingTask

	// ThenDo specifies a task to be run after this one.
	// Returns the new task.
	ThenDo(name string, f Func, opts ...Option) (updated UnderlyingTask, next *Task)

	// WhileDo specifies a task to be run in parallel to this task.
	// Returns the task representing the collection of tasks to be run.
	WhileDo(name string, f Func, opts ...Option) (updated UnderlyingTask)

	// Requires specifies a task to be run before this task.
	// Returns the updated task.
	Requires(*Task) (updated UnderlyingTask)

	// SetContext sets a context to the task.
	// Returns the updated task.
	SetContext(context.Context) (updated UnderlyingTask)

	// Run executes the task.
	Run() error
}
