package alphadns

import (
	"testing"
)

func TestSliceReverse(t *testing.T) {
	s := []string{"hello", "world"}
	rs := reverse(s)
	if rs[0] != "world" || rs[1] != "hello" {
		t.Error("reverse slice broken: ", s, rs)
	}
}
