package param

import (
	"encoding/json"
	"net/http"
)

type ReqGetAuth struct {
	AccessKeyId     string
	AccessKeySecret string
}

func (a ReqGetAuth) Method() string {
	return http.MethodGet
}

func (a ReqGetAuth) Url() string {
	return GetAuthRequestURL
}

func (a ReqGetAuth) Body() []byte {
	return nil
}

func (a ReqGetAuth) Headers() map[string]string {
	headers := make(map[string]string)
	headers["X-Access-KeyId"] = a.AccessKeyId
	headers["X-Access-KeySecret"] = a.AccessKeySecret
	return headers
}

type RspGetAuth struct {
	Resp
	Payload PayloadGetAuth `json:"payload"`
}

func (resp RspGetAuth) String() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

type PayloadGetAuth struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}
