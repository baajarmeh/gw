package gwcall

import (
	assert2 "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExec_Normal(t *testing.T) {
	var f = func() {
		time.Sleep(50 * time.Microsecond)
	}
	err := Call(f, 100*time.Microsecond)
	assert2.True(t, err == nil)
}

func TestExec_Timeout(t *testing.T) {
	var f = func() {
		time.Sleep(500 * time.Microsecond)
	}
	err := Call(f, 100*time.Microsecond)
	assert2.True(t, err != nil)
}

func BenchmarkExec_Normal(b *testing.B) {
	var f = func() {
		time.Sleep(1 * time.Nanosecond)
	}
	for i := 0; i < b.N; i++ {
		_ = Call(f, 5*time.Nanosecond)
	}
}

func BenchmarkExec_Timeout(b *testing.B) {
	var f = func() {
		time.Sleep(5 * time.Nanosecond)
	}
	for i := 0; i < b.N; i++ {
		_ = Call(f, 1*time.Nanosecond)
	}
}
