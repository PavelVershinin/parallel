package parallel

import (
	"sync"
)

// Parallel Паралельно выполнит все функции fns
// Время работы функции Parallel примерно равно времени работы самой долгой из переданных в Parallel фунций
// Для досрочного выхода, закрыть канал cancel
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
