package fsq

func Curry[T, U, V any](f func(T, U) V) func(T) func(U) V {
	return func(t T) func(U) V {
		return func(u U) V {
			return f(t, u)
		}
	}
}

func IgnoreArg[T, U any](f func() (U, error)) func(T) (U, error) {
	return func(_ T) (U, error) {
		return f()
	}
}
