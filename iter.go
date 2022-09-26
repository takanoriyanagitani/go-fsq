package fsq

type Iter[T any] func() (t T, hasValue bool)

func IterReduce[T, U any](i Iter[T], init U, reducer func(state U, t T) U) U {
	var state U = init
	for o, hasValue := i(); hasValue; o, hasValue = i() {
		var t T = o
		state = reducer(state, t)
	}
	return state
}

func IterFromArr[T any](a []T) Iter[T] {
	var ix int = 0
	return func() (t T, hasValue bool) {
		if ix < len(a) {
			t = a[ix]
			ix += 1
			return t, OptHasValue
		}
		return t, OptEmpty
	}
}

func (i Iter[T]) TryForEach(f func(T) error) error {
	return IterReduce(i, nil, func(state error, t T) error {
		return ErrOrElse(state, func() error { return f(t) })
	})
}

func (i Iter[T]) All(f func(T) bool) bool {
	return IterReduce(i, true, func(state bool, t T) bool {
		return state && f(t)
	})
}
