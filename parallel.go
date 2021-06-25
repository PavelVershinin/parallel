package parallel

import (
	"sync"
)

func Parallel(cancel <-chan struct{}, fns ...func()) {
	doneCh := make(chan struct{})

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(fns))
		for _, fn := range fns {
			go func(fn func()) {
				if fn != nil {
					fn()
				}
				wg.Done()
			}(fn)
		}
		wg.Wait()
		doneCh <- struct{}{}
	}()

	go func() {
		<-cancel
		doneCh <- struct{}{}
	}()

	<-doneCh
}
