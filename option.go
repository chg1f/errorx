package errorx

type Option[T comparable] struct {
	empty T
	nan   T
}

func option[T comparable]() Option[T] {
	return Option[T]{}
}

func (ob Option[T]) clone() Option[T] {
	return Option[T]{
		empty: ob.empty,
		nan:   ob.nan,
	}
}

func NaN[T comparable](code T) Option[T] {
	return option[T]().NaN(code)
}

func (ob Option[T]) NaN(code T) Option[T] {
	n := ob.clone()
	n.nan = code
	return n
}

func Empty[T comparable](code T) Option[T] {
	return option[T]().Empty(code)
}

func (ob Option[T]) Empty(code T) Option[T] {
	n := ob.clone()
	n.empty = code
	return n
}

func (ob Option[T]) Code(err error) T {
	ex := Be[T](err)
	if ex == nil {
		return ob.nan
	}
	var empty T
	if ex.code == empty {
		return ob.empty
	}
	return ex.code
}
