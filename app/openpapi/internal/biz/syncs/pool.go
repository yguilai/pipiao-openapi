package syncs

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/threading"
	"sync"
)

type (
	DoSyncFn[T any] func(ctx context.Context, e *T, errch chan<- error)

	SyncPool[T any] struct {
		pool *ants.Pool
	}
)

func NewSyncPool[T any](size int) *SyncPool[T] {
	pool, _ := ants.NewPool(size)
	return &SyncPool[T]{
		pool: pool,
	}
}

func (s *SyncPool[T]) AsyncAll(ctx context.Context, all []T, syncFn DoSyncFn[T]) <-chan error {
	errch := make(chan error, 1)
	threading.GoSafe(func() {
		defer close(errch)
		var wg sync.WaitGroup

		for _, t := range all {
			wg.Add(1)
			err := s.pool.Submit(func() {
				defer wg.Done()
				syncFn(ctx, &t, errch)
			})
			if err != nil {
				errch <- err
			}
		}
		wg.Wait()
	})
	return errch
}

func (s *SyncPool[T]) SyncAll(ctx context.Context, all []T, syncFn DoSyncFn[T]) error {
	errch := s.AsyncAll(ctx, all, syncFn)

	for err := range errch {
		if err != nil {
			return err
		}
	}
	return nil
}
