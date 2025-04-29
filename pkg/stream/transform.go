package stream

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
