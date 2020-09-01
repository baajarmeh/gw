package gw

import (
	"testing"
	"unsafe"
)

func TestName(t *testing.T) {
	var a = make(map[string]bool)
	var p = unsafe.Pointer(&a)
	t.Logf("a addr: %d", p)
	modifyA(t, &a)
}

func modifyA(t *testing.T, a *map[string]bool) {
	for _a, _ := range *a {
		(*a)[_a] = true
	}
	var p = unsafe.Pointer(a)
	t.Logf("a addr: %d", p)
}
