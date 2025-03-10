package pool

import (
	"sync"
)

func StartPool[T any](workerCount int, jobs <-chan Job[T]) <-chan JobResult[T] {
	numberOfWorkers := defaultWorkers
	if workerCount > 0 {
		numberOfWorkers = workerCount
	}

	wg := new(sync.WaitGroup)
	wg.Add(numberOfWorkers)
	jobChannel := make(chan Job[T])
	jobResultChannel := make(chan JobResult[T], 100)

	// Start the workers
	for i := 0; i < numberOfWorkers; i++ {
		go worker(i, wg, jobChannel, jobResultChannel)
	}

	go func() {
		for job := range jobs {
			jobChannel <- job
		}
		close(jobChannel)
		wg.Wait()
		close(jobResultChannel)
	}()

	return jobResultChannel
}
