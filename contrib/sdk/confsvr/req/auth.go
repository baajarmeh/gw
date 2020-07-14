package req

import (
	"encoding/json"
	"net/http"
)

type GetAuthRequest struct {
	AccessKeyId     string
	AccessKeySecret string
}

func (a GetAuthRequest) Method() string {
	return http.MethodGet
}

func (a GetAuthRequest) Url() string {
	return GetAuthRequestURL
}

func (a GetAuthRequest) Body() []byte {
	return nil
}

func (a GetAuthRequest) Headers() map[string]string {
	headers := make(map[string]string)
	headers["X-Access-KeyId"] = a.AccessKeyId
	headers["X-Access-KeySecret"] = a.AccessKeySecret
	return headers
}

type RespGetAuth struct {
	Resp
	Payload PayloadGetAuth `json:"payload"`
}

func (resp RespGetAuth) String() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

type PayloadGetAuth struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}
