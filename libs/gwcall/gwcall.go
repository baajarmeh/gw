package gwcall

import (
	"fmt"
	"time"
)

var ErrorTimeout = fmt.Errorf("timeout")

func Call(f func(), timeout time.Duration) error {
	var done = make(chan bool, 1)
	defer close(done)
	go func(ch chan bool) {
		f()
		ch <- true
	}(done)
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return ErrorTimeout
	}
}
