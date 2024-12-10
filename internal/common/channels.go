package common

import (
	"sync"
)

type Processor[E any, I any, R any] func(E, I) (R, error)

func ParallelExecute[E any, I any, R any](
	inputData I,
	entities []E,
	processor Processor[E, I, R],
) ([]R, error) {
	results := []R{}
	errors := []error{}

	var wg sync.WaitGroup
	resultChan := make(chan R, len(entities))
	errorChan := make(chan error)

	for _, entity := range entities {
		wg.Add(1)

		go func(e interface{}) {
			defer wg.Done()

			result, err := processor(entity, inputData)

			if err != nil {
				errorChan <- err
			} else {
				resultChan <- result
			}
		}(entity)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	for result := range resultChan {
		results = append(results, result)
	}
	for err := range errorChan {
		errors = append(errors, err)
	}

	return results, MergeErrors(errors)

}
