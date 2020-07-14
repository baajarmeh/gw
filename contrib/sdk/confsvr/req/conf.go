package req

import (
	"encoding/json"
	"net/http"
)

type GetConfRequest struct {
	Token       string `json:"token"`
	Namespace   string `json:"ns"`
	Environment string `json:"env"`
}

func (a GetConfRequest) Method() string {
	return http.MethodGet
}

func (a GetConfRequest) Url() string {
	return GetConfRequestURL
}

func (a GetConfRequest) Body() []byte {
	return nil
}

func (a GetConfRequest) Headers() map[string]string {
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
	Version int64 `json:"version"`
}

// ====================================

type CheckConfigVersionRequest struct {
	Token string `json:"token"`
}

func (r CheckConfigVersionRequest) Method() string {
	panic("implement me")
}

func (r CheckConfigVersionRequest) Url() string {
	panic("implement me")
}

func (r CheckConfigVersionRequest) Body() []byte {
	panic("implement me")
}

func (r CheckConfigVersionRequest) Headers() map[string]string {
	panic("implement me")
}

type RespCheckConfigVersion struct {
	Resp
	Payload PayloadCheckConfigVersion `json:"payload"`
}

type PayloadCheckConfigVersion struct {
	Version    int64 `json:"last_version"`
	ModifiedAt int64 `json:"modified_at"`
}
