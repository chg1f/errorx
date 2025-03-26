package errorx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_Code(t *testing.T) {
	a := assert.New(t)

	type code int
	var (
		nan      = code(-2)
		empty    = code(-1)
		nonexist code
		exist    = code(1)
	)
	err0 := build[code]().New("")
	err1 := Code(exist).New("")

	a.Equal(exist, NaN(nan).Code(err1))
	a.Equal(exist, NaN(nan).Empty(empty).Code(err1))
	a.Equal(exist, Empty(empty).Code(err1))
	a.Equal(exist, Empty(empty).NaN(nan).Code(err1))

	a.Equal(nonexist, NaN(nan).Code(err0))
	a.Equal(empty, NaN(nan).Empty(empty).Code(err0))
	a.Equal(empty, Empty(empty).Code(err0))
	a.Equal(empty, Empty(empty).NaN(nan).Code(err0))

	a.Equal(nan, NaN(nan).Code(nil))
	a.Equal(nan, NaN(nan).Empty(empty).Code(nil))
	a.Equal(nonexist, Empty(empty).Code(nil))
	a.Equal(nan, Empty(empty).NaN(nan).Code(nil))
}
