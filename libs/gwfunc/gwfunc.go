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
	go func(e *execution) {
		defer e.release()
		e.f()
		if e.hasTimeout {
			return
		}
		e.hasDone <- true
	}(e)
	select {
	case <-time.After(e.timeout):
		e.hasTimeout = true
		isTimeout <- true
	case <-e.hasDone:
		isTimeout <- false
	}
	return isTimeout
}

func (e *execution) release() {
	close(e.hasDone)
}

// WithTimeout represents a API that can be supports timeout control
// returns false if call f has done (not timeout) else true.
func internalTimeout(f func(), timeout time.Duration) <-chan bool {
	var executor = execution{
		f:       f,
		timeout: timeout,
		hasDone: make(chan bool, 1),
	}
	return executor.exec()
}

// WithTimeout represents a API that can be supports timeout control
// returns false if call f has done (not timeout) else true.
func Timeout(f func(), timeout time.Duration) bool {
	return <-internalTimeout(f, timeout)
}

type TimeSpentResult struct {
	Spent      time.Duration
	isTimeout  bool
	StartedAt  time.Time
	FinishedAt time.Time
}

func (tsr TimeSpentResult) IsTimeout() bool {
	return tsr.isTimeout
}

// TimeSpent represents a API that can be supports record f exec spent time
func TimeSpent(f func(), timeout time.Duration) TimeSpentResult {
	var result TimeSpentResult
	result.StartedAt = time.Now()
	result.isTimeout = Timeout(f, timeout)
	if !result.isTimeout {
		result.FinishedAt = time.Now()
		result.Spent = result.FinishedAt.Sub(result.StartedAt)
	}
	return result
}
