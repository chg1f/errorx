package errorx

import (
	"errors"
	"strconv"
)

func ExampleCompress() {
	var err error
	for i := 0; i < 5; i += 1 {
		err = Compress(err, errors.New(strconv.Itoa(i)))
	}
	if err != nil {
		// error
	}
}
