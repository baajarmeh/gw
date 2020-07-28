package param

import (
	"encoding/json"
	"net/http"
)

type ReqGetConf struct {
	Token       string `json:"token"`
	Namespace   string `json:"ns"`
	Environment string `json:"env"`
}

func (a ReqGetConf) Method() string {
	return http.MethodGet
}

func (a ReqGetConf) Url() string {
	return GetConfRequestURL
}

func (a ReqGetConf) Body() []byte {
	return nil
}

func (a ReqGetConf) Headers() map[string]string {
	return nil
}

type RespGetConf struct {
	Resp
	Payload PayloadGetConf `json:"payload"`
}

func (resp RespGetConf) String() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

type PayloadGetConf struct {
	Data    string `json:"data"`
	Version int64  `json:"version"`
}

// ====================================

type ReqCheckConfigVersion struct {
	Token string `json:"token"`
}

func (r ReqCheckConfigVersion) Method() string {
	panic("implement me")
}

func (r ReqCheckConfigVersion) Url() string {
	panic("implement me")
}

func (r ReqCheckConfigVersion) Body() []byte {
	panic("implement me")
}

func (r ReqCheckConfigVersion) Headers() map[string]string {
	panic("implement me")
}

type RspCheckConfigVersion struct {
	Resp
	Payload PayloadCheckConfigVersion `json:"payload"`
}

type PayloadCheckConfigVersion struct {
	Version    int64 `json:"last_version"`
	ModifiedAt int64 `json:"modified_at"`
}
