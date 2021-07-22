package functions

import (
	"sync"
	"time"
)

func ForWidthBreaks(
	allSteps,
	countStepsBeforeBreak int,
	breakTime time.Duration,
	worker func(wg *sync.WaitGroup, step int) error,
) error {
	for t:=0; t<allSteps; t+=countStepsBeforeBreak {
		countSteps := minInt(
			countStepsBeforeBreak,
			allSteps-t,
		)
		wg := new(sync.WaitGroup)
		wg.Add(countSteps)
		var err error
		for r:=0; r<countSteps; r++ {
			go func(step int){
				workerError := worker(wg, step)
				if workerError != nil {
					err = workerError
				}
			}(t+r)
		}
		wg.Wait()
		<-time.After(breakTime)
		if err != nil {
			return err
		}
	}
	return nil
}