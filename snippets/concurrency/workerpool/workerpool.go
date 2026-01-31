package workerpool

import (
	"context"
	"sync"
)

type Pool[T any] struct {
	workers int
}

func New[T any](workers int) *Pool[T] {
	if workers <= 0 {
		workers = 1
	}
	return &Pool[T]{workers: workers}
}

type job[T any] struct {
	idx int
	val T
}

type result[T any] struct {
	idx int
	val T
	err error
}

// Run executes fn for each input using a fixed number of workers.
//
// It preserves ordering: returned slice aligns with inputs by index.
// On the first error, it cancels remaining work and returns that error.
func (p *Pool[T]) Run(ctx context.Context, inputs []T, fn func(context.Context, T) (T, error)) ([]T, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan job[T])
	results := make(chan result[T])

	var wg sync.WaitGroup
	wg.Add(p.workers)
	for i := 0; i < p.workers; i++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				if ctx.Err() != nil {
					return
				}
				v, err := fn(ctx, j.val)
				select {
				case results <- result[T]{idx: j.idx, val: v, err: err}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for i, in := range inputs {
			select {
			case jobs <- job[T]{idx: i, val: in}:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	outs := make([]T, len(inputs))
	for r := range results {
		if r.err != nil {
			cancel()
			return nil, r.err
		}
		outs[r.idx] = r.val
	}
	return outs, nil
}
