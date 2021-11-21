package parallel

import (
	"context"
	"sync"
)

// Parallel Паралельно выполнит все функции fns
// Время работы функции Parallel примерно равно времени работы самой долгой из переданных в Parallel фунций
// По завершении вернёт слайс ошибок
func Parallel(ctx context.Context, fns ...func() error) []error {
	errors := make([]error, len(fns))
	doneCh := make(chan struct{})

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(fns))
		for i, fn := range fns {
			go func(fNum int, fn func() error) {
				if fn != nil {
					errors[fNum] = fn()
				}
				wg.Done()
			}(i, fn)
		}
		wg.Wait()
		doneCh <- struct{}{}
	}()

	go func() {
		<-ctx.Done()
		doneCh <- struct{}{}
	}()

	<-doneCh
	return errors
}
