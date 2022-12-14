package fsq

func ErrOrElse(e error, ef func() error) error {
	if nil != e {
		return e
	}
	return ef()
}

func Err1st(ef []func() error) error {
	var ei Iter[func() error] = IterFromArr(ef)
	return IterReduce(ei, nil, ErrOrElse)
}

func ErrFuncGen[T, U any](f func(T) U) func(T) (U, error) {
	return func(t T) (U, error) {
		return f(t), nil
	}
}

func ErrOnly[T, U any](f func(T) (U, error)) func(T) error {
	return func(t T) error {
		_, e := f(t)
		return e
	}
}

func ErrFromBool[T any](ok bool, okf func() T, ngf func() error) (t T, e error) {
	if ok {
		return okf(), nil
	}
	return t, ngf()
}

func ErrUnwrapOrElse[T, U any](f func(T) (U, error), g func(error) U) func(T) U {
	return func(t T) U {
		u, e := f(t)
		if nil != e {
			return g(e)
		}
		return u
	}
}

func ErrTryForEach[T any](t T, e error, f func(T) error) error {
	if nil != e {
		return e
	}
	return f(t)
}
