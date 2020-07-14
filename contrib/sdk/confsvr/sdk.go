package confsvr

import (
	"fmt"
	"github.com/oceanho/gw/contrib/sdk/confsvr/req"
	"sync"
	"time"
)

var (
	sdk  *SdkSetting
	once sync.Once
)

type SdkSetting struct {
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

func NewSdkSetting(ak, aks, ns, env string) *SdkSetting {
	return &SdkSetting{
		AccessKeyId:     ak,
		AccessKeySecret: aks,
		Namespace:       ns,
		Environment:     env,
	}
}

func Initial(conf *SdkSetting) {
	once.Do(func() {
		sdk = conf
	})
}

func SyncConf() ([]byte, error) {
	updateToken()
	updateConf()
	return []byte(sdk.currentConfigData), nil
}

func StartWatcher(opts *Option, shutdownSignal chan struct{}) {
	if sdk == nil {
		panic("should be call confsvr.Initial(...) At first.")
	}
	// Initial vars.
	sdk.client = NewClient(opts)

	// 1. Got token At first.
	updateToken()

	// 2. Got configuration.
	updateConf()

	go watchWorker(opts)
	go func() {
		// wait for showdown server Signal.
		<-shutdownSignal
		sdk.shutdown = true
	}()
}

func updateToken() {
	reqobj := req.GetAuthRequest{
		AccessKeyId:     sdk.AccessKeyId,
		AccessKeySecret: sdk.AccessKeySecret,
	}
	respobj := &req.RespGetAuth{}
	code, err := sdk.client.Do(reqobj, respobj)
	if err != nil {
		fmt.Printf("[configsvr-worker] - [WARNING, updateToken fail, code:%d, err: %v", code, err)
		return
	}
	sdk.token = respobj.Payload.Token
	sdk.tokenExpiredAt = respobj.Payload.ExpiredAt
}

func hasNewVersion() bool {
	reqobj := req.CheckConfigVersionRequest{
		Token: sdk.token,
	}
	respobj := &req.RespCheckConfigVersion{}
	code, err := sdk.client.Do(reqobj, respobj)
	if err != nil {
		fmt.Printf("[configsvr-worker] - [WARNING, hasNewVersion fail, code:%d, err: %v", code, err)
	}
	return respobj.Payload.Version > sdk.currentConfigVersion
}

func updateConf() {
	reqobj := req.GetConfRequest{
		Token: sdk.token,
	}
	respobj := &req.RespGetConf{}
	code, err := sdk.client.Do(reqobj, respobj)
	if err != nil {
		fmt.Printf("[configsvr-worker] - [WARNING, updateConf fail, code:%d, err: %v", code, err)
		return
	}
	sdk.currentConfigData = respobj.Payload.Data
	sdk.currentConfigVersion = respobj.Payload.Version
}

func watchWorker(opts *Option) {
	interval := time.Duration(opts.PullInterval) * time.Second
	for {
		select {
		case <-time.After(interval):
			if hasNewVersion() {
				updateConf()
			}
		}
		if sdk.shutdown {
			break
		}
	}
}
