package errorx

type CodeOption[T comparable] struct {
	Empty T
	NaN   T
}

func NaN[T comparable](code T) func(*CodeOption[T]) {
	return func(opt *CodeOption[T]) {
		opt.NaN = code
	}
}

func Empty[T comparable](code T) func(*CodeOption[T]) {
	return func(opt *CodeOption[T]) {
		opt.Empty = code
	}
}

func CodeOf[T comparable](err error, opts ...func(*CodeOption[T])) T {
	var opt CodeOption[T]
	for i := range opts {
		opts[i](&opt)
	}
	ex := Be[T](err)
	if ex == nil {
		return opt.NaN
	}
	var empty T
	code := ex.Code()
	if ex.code == empty {
		return opt.Empty
	}
	return code
}
