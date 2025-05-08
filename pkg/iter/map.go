package iter

func Map[T any, U any](s []T, mapFunc func(T) U) []U {
	result := make([]U, 0, len(s))
	for _, v := range s {
		result = append(result, mapFunc(v))
	}
	return result
}

func Reduce[T any, A any](s []T, reduceFunc func(A, T) A, initial A) A {
	acc := initial
	for _, v := range s {
		acc = reduceFunc(acc, v)
	}
	return acc
}
