package collect

import "github.com/zeromicro/go-zero/core/threading"

type (
	Slice[T any] []T

	Map[T comparable, U any] map[T]U

	Iterator[T any] interface {
		It() <-chan T
	}
)

func (m Map[T, U]) It() <-chan U {
	r := make(chan U, 1)
	threading.GoSafe(func() {
		for _, t := range m {
			r <- t
		}
	})
	return r
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
