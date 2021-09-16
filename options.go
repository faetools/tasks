package tasks

import "context"

// Option applies some setting to a task.
type Option func(UnderlyingTask) UnderlyingTask

// OptionRequires makes sure another task is executed beforehand.
func OptionRequires(required *Task) Option {
	return func(t UnderlyingTask) UnderlyingTask {
		return t.Requires(required)
	}
}

// OptionContext provides the tasks with a context.
func OptionContext(ctx context.Context) Option {
	return func(t UnderlyingTask) UnderlyingTask {
		return t.SetContext(ctx)
	}
}

// // OptionTimeout adds a timeout to the task.
// // The timer starts as soon as the task is about to be started.
// func OptionTimeout(timeout time.Duration) Option {
// 	return func(t Task) Task {
// 		switch v := t.(type) {
// 		case *task:
// 			v.timeout = timeout
// 			return v
// 		default:
// 			panic(fmt.Sprintf("unknown task type %T", v))
// 		}
// 	}
// }

// // OptionRetries adds the option to retry the task a couple of times.
// // A negative number means the task is not even done once.
// func OptionRetries(retries int) Option {
// 	if retries < 0 {
// 		retries = -1
// 	}
// 	return func(t Task) Task {
// 		switch v := t.(type) {
// 		case *task:
// 			v.retries = retries
// 			return v
// 		default:
// 			panic(fmt.Sprintf("unknown task type %T", v))
// 		}
// 	}
// }

// // OptionLogger adds the option to log certain events.
// func OptionLogger(logger io.Writer) Option {
// 	return func(t Task) Task {
// 		switch v := t.(type) {
// 		case *task:
// 			v.logger = logger
// 			return v
// 		default:
// 			panic(fmt.Sprintf("unknown task type %T", v))
// 		}
// 	}
// }
