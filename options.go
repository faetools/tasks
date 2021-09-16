package tasks

import (
	"context"
	"io"
	"time"
)

// Option applies some setting to a task.
type Option func(*Task)

// OptionRequires makes sure another task is executed beforehand.
func OptionRequires(required *Task) Option {
	return func(t *Task) { t.Requires(required) }
}

// OptionContext provides the tasks with a context.
func OptionContext(ctx context.Context) Option {
	return func(t *Task) { t.SetContext(ctx) }
}

// OptionTimeout adds a timeout to the task.
// The timer starts as soon as the task is about to be started.
func OptionTimeout(timeout time.Duration) Option {
	return func(t *Task) { t.timeout = timeout }
}

// OptionRetries adds the option to retry the task a couple of times.
// A negative number means the task is not even done once.
func OptionRetries(retries int) Option {
	return func(t *Task) { t.retries = retries }
}

// OptionLogger adds the option to log certain events.
func OptionLogger(logger io.Writer) Option {
	return func(t *Task) { t.logger = logger }
}
