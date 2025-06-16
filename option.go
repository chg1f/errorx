package errorx

type Option[T comparable] struct {
	Empty T
	NaN   T
}

func NaN[T comparable](code T) Option[T] {
	return Option[T]{
		NaN: code,
	}
}

func Empty[T comparable](code T) Option[T] {
	return Option[T]{
		Empty: code,
	}
}

func CodeOf[T comparable](err error, opts ...Option[T]) T {
	if len(opts) > 1 {
		panic("only one option is allowed")
	}
	var empty T
	ex := Be[T](err)
	if ex == nil {
		if len(opts) == 1 {
			return opts[0].NaN
		}
		return empty
	}
	code := ex.Code()
	if ex.code == empty {
		if len(opts) == 1 {
			return opts[0].Empty
		}
		return empty
	}
	return code
}
