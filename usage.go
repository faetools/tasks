package tasks

// FirstDo specifies a first task to be done.
func FirstDo(name string, f Func, opts ...Option) *Task {
	return newTask(name, f, opts...)
}

// ThenDo specifies a task to be run after this one.
// Returns the new task.
func (t *Task) ThenDo(name string, f Func, opts ...Option) (next *Task) {
	next = newTask(name, f, opts...)
	t.then = append(t.then, next)

	next.Requires(t)
	return next
}

// WhileDo specifies a task to be run in parallel to this task.
func (t *Task) WhileDo(name string, f Func, opts ...Option) *Task {
	t.while = append(t.while, newTask(name, f, opts...))
	return t
}
