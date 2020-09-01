package gwfunc

import (
	//"sync"
	"time"
)

func init() {

}

//var pool = sync.Pool {
//	New: func() interface{}{
//		var executor = &execution {
//			hasTimeout: false,
//			hasDone: make(chan bool, 1),
//		}
//		return executor
//	},
//}

type execution struct {
	f          func()
	state      uint8
	hasDone    chan bool
	hasTimeout bool
	timeout    time.Duration
}

func (e *execution) exec() <-chan bool {
	var isTimeout = make(chan bool, 1)
	go func() {
		e.f()
		if e.hasTimeout {
			return
		}
		e.hasDone <- true
	}()
	select {
	case <-time.After(e.timeout):
		e.hasTimeout = true
		isTimeout <- true
	case <-e.hasDone:
		isTimeout <- false
	}
	e.release()
	return isTimeout
}

func (e *execution) release() {
	close(e.hasDone)
}

// WithTimeout represents a API that can be supports timeout control
// returns true if call f has done (not timeout) else false.
func Timeout(f func(), timeout time.Duration) bool {
	var executor = execution{
		f:       f,
		timeout: timeout,
		hasDone: make(chan bool, 1),
	}
	return <- executor.exec()
}

//func WaitAll(timeout time.Duration, funcList ...func()) bool {
//
//}
//
//func WaitOne(timeout time.Duration, funcList ...func()) bool {
//}
