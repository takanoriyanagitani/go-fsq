package fsq

func ComposeErr[T, U, V any](f func(T) (U, error), g func(U) (V, error)) func(T) (V, error) {
	return func(t T) (v V, e error) {
		u, e := f(t)
		if nil != e {
			return v, e
		}
		return g(u)
	}
}

func Compose[T, U, V any](f func(T) U, g func(U) V) func(T) V {
	return func(t T) V {
		var h func(T) (V, error) = ComposeErr(
			ErrFuncGen(f),
			ErrFuncGen(g),
		)
		v, _ := h(t)
		return v
	}
}
