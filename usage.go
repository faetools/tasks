package tasks

// FirstDo specifies a first task to be done.
func FirstDo(name string, f Func, opts ...Option) *Task {
	return newTask(name, f, opts...)
}

// ThenDo specifies a task to be run after this one.
// Returns the new task.
func (t *Task) ThenDo(name string, f Func, opts ...Option) *Task {
	updated, next := t.t.ThenDo(name, f, opts...)
	t.t = updated
	return next
}

// WhileDo specifies a task to be run in parallel to this task.
// Returns the task representing the collection of tasks to be run.
func (t *Task) WhileDo(name string, f Func, opts ...Option) *Task {
	t.t = t.t.WhileDo(name, f, opts...)
	return t
}
