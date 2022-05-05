package errorx

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShrink(t *testing.T) {
	es := []error{os.ErrNotExist, os.ErrExist}
	temp := make([]string, 0, len(es))
	for _, e := range es {
		temp = append(temp, e.Error())
	}
	text := strings.Join(temp, "; ")
	ce := Compress(es...)

	assert.Equal(t, text, ce.Error())
	assert.Equal(t, text, Shrink(ce).Error())

	fmt.Println(ce.Error())
}

func TestCompress(t *testing.T) {
	errPath0 := &os.PathError{Path: strconv.FormatUint(rand.Uint64(), 16), Err: os.ErrExist}
	errPath1 := &os.PathError{Path: strconv.FormatUint(rand.Uint64(), 16), Err: os.ErrExist}
	ce := Compress(os.ErrExist, errPath0, errPath1)

	assert.True(t, errors.Is(ce, os.ErrExist))
	assert.True(t, errors.Is(ce, errPath0))
	assert.True(t, errors.Is(ce, errPath1))
	assert.False(t, errors.Is(ce, os.ErrNotExist))

	{
		var e *os.PathError
		assert.True(t, errors.As(ce, &e))
		assert.Equal(t, errPath0.Path, e.Path)
		assert.NotEqual(t, errPath1.Path, e.Path)
	}
	{
		var e *CompressError
		assert.True(t, errors.As(ce, &e))
		var e1 *os.PathError
		assert.True(t, errors.As(e.Errors[2], &e1))
		assert.Equal(t, errPath1.Path, e1.Path)
	}
}

func TestCompressNil(t *testing.T) {
	assert.Nil(t, Compress(nil, nil))
	assert.Len(t, Compress(nil, os.ErrClosed).(*CompressError).Errors, 1)
}

func ExampleShrink() {
	fmt.Println(Shrink(os.ErrNotExist, os.ErrExist))
	// output: file does not exist; file already exists
}

func ExampleCompress() {
	ce := Compress(os.ErrNotExist, os.ErrExist)
	fmt.Println(ce)
	// output: file does not exist; file already exists
}
