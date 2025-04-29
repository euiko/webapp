package stream

type (
	Continuation[T any] func(T)
	Stream[T any]       struct {
		next func(Continuation[T])
	}

	MapFunc[T any, E any] func(T) E
	FilterFunc[T any]     func(T) bool
)

func SliceStream[T any](s []T) Stream[T] {
	return Stream[T]{
		next: func(next Continuation[T]) {
			for _, v := range s {
				next(v)
			}
		},
	}
}

func ChanStream[T any](c <-chan T) Stream[T] {
	return Stream[T]{
		next: func(next Continuation[T]) {
			for v := range c {
				next(v)
			}
		},
	}
}
