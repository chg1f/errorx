package errorx

import (
	"errors"
	"strings"
)

// FIXME: WeakError?CompressError?

func Compress(es ...error) error {
	t := make([]string, 0, len(es))
	for _, e := range es {
		if e != nil {
			t = append(t, e.Error())
		}
	}
	if len(t) == 0 {
		return nil
	}
	return errors.New(strings.Join(t, "; "))
}
