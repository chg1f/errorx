package errorx

import (
	"errors"
	"testing"
)

func TestCompressError(t *testing.T) {
	es := new(CompressError)
	if es.Error() != "" {
		t.Error(es.Error())
	}
	if es.x != nil {
		t.Error(es.x)
	}
	es.Append(errors.New("A"))
	if es.Error() != "A" {
		t.Error(es.Error())
	}
	if len(es.x) != 1 {
		t.Error(es.x)
	}
	es.Append(errors.New("B"))
	if es.Error() != "A;B" {
		t.Error(es.Error())
	}
	if len(es.x) != 2 {
		t.Error(es.x)
	}
	e := es.Shrink()
	if e.Error() != "A;B" {
		t.Error(es.Error())
	}
	if es.Error() != "A;B" {
		t.Error(es.Error())
	}
	if len(es.x) != 2 {
		t.Error(es.x)
	}
}
