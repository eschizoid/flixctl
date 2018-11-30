package worker

import (
	"context"
)

type TaskFunction func() interface{}

func PerformTasks(ctx context.Context, tasks []TaskFunction) {
	for _, task := range tasks {
		//resultChannel := newWorker(ctx, task)
		newWorker(ctx, task)
	}
}

func newWorker(ctx context.Context, task TaskFunction) chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		select {
		case <-ctx.Done():
			// Received a signal to abandon further processing
			return
		case out <- task():
			// Got some result
		}
	}()
	return out
}
