package gw

import (
	"testing"
)

func TestName(t *testing.T) {
	var a = make(map[string]int)
	modifyA(t, a)
	for k, v := range a {
		t.Logf("k=v: %s,%d", k, v)
	}
}

func modifyA(t *testing.T, a map[string]int) {
	a["a"] = 100
}
