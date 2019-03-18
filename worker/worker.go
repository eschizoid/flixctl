package worker

import (
	"context"
	"sync"
)

type TaskFunction func() interface{}

func PerformTasks(ctx context.Context, tasks []TaskFunction) chan interface{} {

	// Create a worker for each incoming task
	workers := make([]chan interface{}, 0, len(tasks))

	for _, task := range tasks {
		resultChannel := newWorker(ctx, task)
		workers = append(workers, resultChannel)
	}

	// Merge results from all workers
	out := merge(ctx, workers)
	return out
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

func merge(ctx context.Context, workers []chan interface{}) chan interface{} {
	// Merged channel with results
	out := make(chan interface{})

	// Synchronization over channels: do not close "out" before all tasks are completed
	var wg sync.WaitGroup

	// Define function which waits the result from worker channel
	// and sends this result to the merged channel.
	// Then it decreases the counter of running tasks via wg.Done().
	output := func(c <-chan interface{}) {
		defer wg.Done()
		for result := range c {
			select {
			case <-ctx.Done():
				// Received a signal to abandon further processing
				return
			case out <- result:
				// some message or nothing
			}
		}
	}

	wg.Add(len(workers))
	for _, workerChannel := range workers {
		go output(workerChannel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
