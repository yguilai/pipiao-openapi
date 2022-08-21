package syncs

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"github.com/yguilai/pipiao-openapi/common/collect"
	"github.com/zeromicro/go-zero/core/threading"
	"sync"
)

type (
	DoSyncFn[T any] func(ctx context.Context, e *T, errch chan<- error)

	SyncPool[T collect.Iterator[U], U any] struct {
		pool *ants.Pool
	}
)

func NewSyncPool[T collect.Iterator[U], U any](size int) *SyncPool[T, U] {
	pool, _ := ants.NewPool(size)
	return &SyncPool[T, U]{
		pool: pool,
	}
}

func (s *SyncPool[T, U]) AsyncAll(ctx context.Context, all T, syncFn DoSyncFn[U]) <-chan error {
	errch := make(chan error, 1)
	threading.GoSafe(func() {
		defer close(errch)
		var wg sync.WaitGroup

		c := any(all).(collect.Iterator[U])

		for u := range c.It() {
			wg.Add(1)
			err := s.pool.Submit(func() {
				defer wg.Done()
				syncFn(ctx, &u, errch)
			})
			if err != nil {
				errch <- err
			}
		}
		wg.Wait()
	})
	return errch
}

func (s *SyncPool[T, U]) SyncAll(ctx context.Context, all T, syncFn DoSyncFn[U]) error {
	errch := s.AsyncAll(ctx, all, syncFn)

	for err := range errch {
		if err != nil {
			return err
		}
	}
	return nil
}
