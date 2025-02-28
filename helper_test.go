package errorx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	errorCode := "specified"

	err := New("...")
	assert.Error(t, err)
	assert.True(t, In(err, Unspecified))
	assert.False(t, In(err, errorCode))

	err = Code(errorCode).Wrap(err)
	assert.Error(t, err)
	assert.True(t, In(err, Unspecified))
	assert.True(t, In(err, errorCode))
}

func ExampleIn() {
	err := New("...")
	if In(err, Unspecified) {
		// do something
	}
}

func TestBe(t *testing.T) {
	assert.Nil(t, Be[struct{}](nil))
	assert.Nil(t, Be[int](nil))
	assert.Nil(t, Be[string](nil))

	err := New("...")
	assert.Equal(t, Unspecified, Be[struct{}](err).Code())
	assert.Equal(t, 0, Be[int](err).Code())
	assert.Equal(t, "", Be[string](err).Code())
}

func ExampleBe() {
	ex := Be[struct{}](New("..."))
	switch ex.Code() {
	case Unspecified:
		// do something
	default:
		// do something
	}
}
