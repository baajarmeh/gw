package confsvr

import (
	"fmt"
	"github.com/oceanho/gw/sdk/confsvr/param"
	"sync"
	"time"
)

var st *state

type state struct {
	sync.Once
	state *setting
}

func init() {
	st = &state{}
}

type setting struct {
	AccessKeyId          string
	AccessKeySecret      string
	OnChangedCallback    func(data []byte)
	Namespace            string
	Environment          string
	token                string
	client               *Client
	tokenExpiredAt       int64
	shutdown             bool
	currentConfigVersion int64
	currentConfigData    string
	locker               sync.Mutex
}

func Initial(ak, aks, ns, env string, onDataChangedCallback func(data []byte)) {
	st.Once.Do(func() {
		st.state = &setting{
			AccessKeyId:       ak,
			AccessKeySecret:   aks,
			Namespace:         ns,
			Environment:       env,
			OnChangedCallback: onDataChangedCallback,
		}
	})
}

func Sync() ([]byte, error) {
	if st.state == nil {
		panic("should be call confsvr.Initial(...) At first.")
	}
	updateToken()
	updateConf()
	return []byte(st.state.currentConfigData), nil
}

func StartWatcher(opts *Option, shutdownSignal chan struct{}) {
	if st.state == nil {
		panic("should be call confsvr.Initial(...) At first.")
	}
	// Initial vars.
	st.state.client = NewClient(opts)

	// 1. Got token At first.
	updateToken()

	// 2. Got configuration.
	updateConf()

	go watchWorker(opts)
	go func() {
		// wait for showdown server Signal.
		<-shutdownSignal
		st.state.shutdown = true
	}()
}

func updateToken() {
	req := param.ReqGetAuth{
		AccessKeyId:     st.state.AccessKeyId,
		AccessKeySecret: st.state.AccessKeySecret,
	}
	resp := &param.RspGetAuth{}
	code, err := st.state.client.Do(req, resp)
	if err != nil {
		fmt.Printf("[configsvr-worker] - [WARNING, updateToken fail, code:%d, err: %v", code, err)
		return
	}
	st.state.token = resp.Payload.Token
	st.state.tokenExpiredAt = resp.Payload.ExpiredAt
}

func hasNewVersion() bool {
	req := param.ReqCheckConfigVersion{
		Token: st.state.token,
	}
	resp := &param.RspCheckConfigVersion{}
	code, err := st.state.client.Do(req, resp)
	if err != nil {
		fmt.Printf("[configsvr-worker] - [WARNING, hasNewVersion fail, code:%d, err: %v", code, err)
	}
	return resp.Payload.Version > st.state.currentConfigVersion
}

func updateConf() {
	req := param.ReqGetConf{
		Token: st.state.token,
	}
	resp := &param.RespGetConf{}
	code, err := st.state.client.Do(req, resp)
	if err != nil {
		fmt.Printf("[configsvr-worker] - [WARNING, updateConf fail, code:%d, err: %v", code, err)
		return
	}
	st.state.currentConfigData = resp.Payload.Data
	st.state.currentConfigVersion = resp.Payload.Version
}

func watchWorker(opts *Option) {
	interval := time.Duration(opts.QueryStateInterval) * time.Second
	for {
		select {
		case <-time.After(interval):
			if hasNewVersion() {
				updateConf()
			}
		}
		if st.state.shutdown {
			break
		}
	}
}
