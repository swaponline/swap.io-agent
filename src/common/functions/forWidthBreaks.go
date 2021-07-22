package functions

import (
	"sync"
	"time"
)

func forWidthBreaks(
	allSteps,
	countStepsBeforeBreak int,
	breakTime time.Duration,
	worker func(wg *sync.WaitGroup, step int),
) {
	for t:=0; t<allSteps; t+=countStepsBeforeBreak {
		countSteps := minInt(
			countStepsBeforeBreak,
			allSteps-t,
		)
		wg := new(sync.WaitGroup)
		wg.Add(countSteps)
		for r:=0; r<countSteps; r++ {
			worker(wg, t+r)
		}
		wg.Wait()
		<-time.After(breakTime)
	}
}