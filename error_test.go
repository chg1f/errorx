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
	var (
		err0 = &os.PathError{Path: strconv.FormatUint(rand.Uint64(), 16), Err: os.ErrExist}
		err1 = errors.New(strconv.FormatUint(rand.Uint64(), 16))
		err2 = errors.New(strconv.FormatUint(rand.Uint64(), 16))
	)
	assert.Nil(t, Shrink(nil))
	assert.Equal(t, strings.Join([]string{err0.Error(), err1.Error(), err2.Error()}, ErrorSeparator), Shrink(err0, err1, err2).Error())
}
func TestCompress(t *testing.T) {
	var (
		err0 = &os.PathError{Path: strconv.FormatUint(rand.Uint64(), 16), Err: os.ErrExist}
		err1 = errors.New(strconv.FormatUint(rand.Uint64(), 16))
		err2 = errors.New(strconv.FormatUint(rand.Uint64(), 16))
	)
	ce0 := Compress(nil)
	assert.Nil(t, ce0)
	{
		ce := new(CompressError)
		assert.False(t, errors.As(ce0, &ce))
	}

	ce1 := Compress(err0, err1)
	assert.NotNil(t, ce1)
	assert.True(t, errors.Is(ce1, err0))
	assert.True(t, errors.Is(ce1, err1))
	assert.False(t, errors.Is(ce1, err2))
	{
		ce := new(CompressError)
		assert.True(t, errors.As(ce1, &ce))
		assert.NotNil(t, ce)
		assert.Len(t, ce.x, 2)
		assert.Len(t, ce.x, ce.Len())
		assert.Equal(t, ce.x[0], err0)
		assert.Equal(t, ce.x[1], err1)
		ce.ForEach(func(_ int, err error) bool {
			assert.NotNil(t, err)
			assert.False(t, errors.As(err, new(*CompressError)))
			return true
		})
	}

	ce2 := Compress(ce0, ce1, err2)
	assert.True(t, errors.Is(ce2, err0))
	assert.True(t, errors.Is(ce2, err1))
	assert.True(t, errors.Is(ce2, err2))
	{
		ce := new(CompressError)
		assert.True(t, errors.As(ce2, &ce))
		assert.NotNil(t, ce)
		assert.Len(t, ce.x, 3)
		assert.Len(t, ce.x, ce.Len())
		assert.Equal(t, ce.x[0], err0)
		assert.Equal(t, ce.x[1], err1)
		assert.Equal(t, ce.x[2], err2)
		ce.ForEach(func(_ int, err error) bool {
			assert.NotNil(t, err)
			assert.False(t, errors.As(err, new(*CompressError)))
			return true
		})
	}
}
func TestForEach(t *testing.T) {
	var (
		err0 = &os.PathError{Path: strconv.FormatUint(rand.Uint64(), 16), Err: os.ErrExist}
		err1 = errors.New(strconv.FormatUint(rand.Uint64(), 16))
		err2 = errors.New(strconv.FormatUint(rand.Uint64(), 16))
	)
	ce := Compress(err0, err1, err2)
	ForEach(ce, func(index int, err error) bool {
		switch index {
		case 0:
			assert.True(t, errors.Is(err, err0))
		case 1:
			assert.True(t, errors.Is(err, err1))
		case 2:
			assert.True(t, errors.Is(err, err2))
		}
		return true
	})
}

func ExampleShrink() {
	var (
		err0 = errors.New("Hello")
		err1 = errors.New("ErrorX")
	)
	fmt.Println(Shrink(err0, nil, err1))
	// output:
	// Hello; ErrorX
}
func ExampleCompress() {
	var (
		err0 = errors.New("Hello")
		err1 = errors.New("ErrorX")
	)
	fmt.Println(Compress(err0, nil, err1).Error())
	// output:
	// Hello; ErrorX
}
func ExampleForEach() {
	var (
		err0 = errors.New("Hello")
		err1 = errors.New("ErrorX")
	)
	ForEach(Compress(err0, nil, err1), func(index int, err error) bool {
		fmt.Println(index, err.Error())
		return true
	})
	// output:
	// 0 Hello
	// 1 ErrorX
}
