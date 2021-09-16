package tasks

import (
	"context"
	"io"
	"sync"
	"time"
)

type meta struct {
	ctx  context.Context
	name string

	requires tasks // what to do beforehand
	then     tasks // what to do afterwards
	timeout  time.Duration
	retries  int
	logger   io.Writer

	done bool
	once sync.Once
	err  error

	container *Task
}

func (m *meta) Apply(o Option) UnderlyingTask { return o(m.container.t) }

func (m *meta) SetContext(ctx context.Context) UnderlyingTask {
	m.ctx = ctx
	return m.container.t
}

func (m *meta) ThenDo(name string, f Func, opts ...Option) (UnderlyingTask, *Task) {
	next := newTask(name, f, opts...)
	m.then = append(m.then, next)
	next.Requires(m.container)

	return m.container.t, next
}

func (m *meta) Requires(required *Task) UnderlyingTask {
	m.requires = append(m.requires, required)
	return m.container.t
}
