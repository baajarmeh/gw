package gwcall

import (
	"fmt"
	"time"
)

var ErrorTimeout = fmt.Errorf("timeout")

func Call(f func(), timeout time.Duration) error {
	var ch = make(chan bool, 1)
	defer close(ch)
	go func() {
		f()
		ch <- true
	}()
	select {
	case _, _ = <-ch:
		return nil
	case <-time.After(timeout):
		return ErrorTimeout
	}
}
