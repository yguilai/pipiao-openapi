package maps

type Entry[T comparable, U any] struct {
	Key   T
	Value U
}

func Entries[T comparable, U any](m map[T]U) []Entry[T, U] {
	es := make([]Entry[T, U], len(m))
	for k, v := range m {
		es = append(es, Entry[T, U]{
			Key:   k,
			Value: v,
		})
	}
	return es
}
