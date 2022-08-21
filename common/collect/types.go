package collect

import "github.com/zeromicro/go-zero/core/threading"

type (
	Slice[T any] []T

	Map[T comparable, U any] map[T]U

	Entry[T comparable, U any] struct {
		Key   T
		Value U
	}

	Iterator[T any] interface {
		It() <-chan T
	}
)

func (m Map[T, U]) Entries() Slice[Entry[T, U]] {
	es := make(Slice[Entry[T, U]], len(m))
	for k, v := range m {
		es = append(es, Entry[T, U]{
			Key:   k,
			Value: v,
		})
	}
	return es
}

func (s Slice[T]) It() <-chan T {
	r := make(chan T, 1)
	threading.GoSafe(func() {
		for _, t := range s {
			r <- t
		}
	})
	return r
}
