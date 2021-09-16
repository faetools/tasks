package tasks

import (
	"context"
	"time"
)

// h/t https://www.opsdash.com/blog/job-queues-in-go.html

// DoWithContext executes a function with the given context.
// If the context is cancelled or timeout happens, the respective
// error is returned.
func DoWithContext(ctx context.Context, f Func) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	ch := make(chan struct{})

	var err error
	go func() {
		err = f()
		close(ch)
	}()

	if dl, ok := ctx.Deadline(); ok {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			return err
		case <-time.After(time.Until(dl)):
			return context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return err
	}
}
