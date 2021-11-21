package parallel_test

import (
	"context"
	"math"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PavelVershinin/parallel"
	"github.com/stretchr/testify/require"
)

func makeTestFuncs(count int, res *int32) (fns []func() error, sumSleep, maxSleep time.Duration) {
	for i := 0; i < count; i++ {
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
		maxSleep = time.Duration(math.Max(float64(maxSleep.Milliseconds()), float64(taskSleep.Milliseconds()))) * time.Millisecond
		sumSleep += taskSleep
		fns = append(fns, func() error {
			<-time.After(taskSleep)
			atomic.AddInt32(res, 1)
			return  nil
		})
	}

	return
}

func TestParallel(t *testing.T) {
	var res int32
	fns, sumSleep, maxSleep := makeTestFuncs(15, &res)

	t.Run("the parallel function must be completed before the sum of the durations of all functions from the array and all functions from the array must have time to be executed", func(t *testing.T) {
		res = 0

		ctx := context.Background()
		start := time.Now()
		parallel.Parallel(ctx, fns...)
		elapsed := time.Since(start)

		require.Less(t, elapsed, sumSleep)
		require.Equal(t, int32(15), atomic.LoadInt32(&res))
	})

	t.Run("the parallel function must be completed earlier than the longest function from the array and not all functions from the array should have time to be executed", func(t *testing.T) {
		res = 0
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(maxSleep / 2))

		start := time.Now()
		parallel.Parallel(ctx, fns...)
		elapsed := time.Since(start)

		require.Less(t, elapsed, maxSleep)
		require.Less(t, atomic.LoadInt32(&res), int32(15))
	})

	t.Run("the number of errors returned must be equal to the number of functions passed", func(t *testing.T) {
		errs := parallel.Parallel(context.Background(), fns...)
		require.Equal(t, len(fns), len(errs))
	})
}
