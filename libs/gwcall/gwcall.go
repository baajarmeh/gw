package gwcall

import (
	"time"
)

// Call represents a API that can be supports timeout control
// returns true if call f has done (not timeout) else false.
func Call(f func(), timeout time.Duration) bool {
	var done = make(chan bool, 1)
	var cancel = make(chan bool, 1)
	defer close(cancel)
	go func() {
		defer close(done)
		for {
			select {
			case _, _ = <-cancel:
				return
			default:
				f()
				done <- true
			}
		}
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		cancel <- true
		return false
	}
}
