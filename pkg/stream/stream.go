package stream

type (
	Continuation[T any] func(T)
	Stream[T any]       struct {
		next func(Continuation[T])
	}

	MapFunc[T any, E any] func(T) E
	FilterFunc[T any]     func(T) bool
)

func Map[T any, E any](s Stream[T], mapFunc MapFunc[T, E]) Stream[E] {
	return Stream[E]{
		next: func(next Continuation[E]) {
			s.next(func(t T) {
				e := mapFunc(t)
				next(e)
			})
		},
	}
}

func Filter[T any](s Stream[T], filterFunc FilterFunc[T]) Stream[T] {
	return Stream[T]{
		next: func(next Continuation[T]) {
			s.next(func(t T) {
				ok := filterFunc(t)
				if ok {
					next(t)
				}
			})
		},
	}
}

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
