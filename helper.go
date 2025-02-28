package errorx

func In[T comparable](err error, code T, from ...string) bool {
	if len(from) > 1 {
		panic("errx.In: too many arguments")
	}
	for {
		if err == nil {
			return false
		}
		if x, ok := err.(interface{ In(T, ...string) bool }); ok {
			return x.In(code)
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if In(err, code) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

func Be[T comparable](err error) *Error[T] {
	if err == nil {
		return nil
	}
	ex, ok := err.(*Error[T])
	if !ok {
		return build[T]().Wrap(err).(*Error[T])
	}
	return ex
}

func Stack(err error) []Frame {
	if ex := Be[struct{}](err); ex != nil {
		return ex.Stack()
	}
	return []Frame{}
}
