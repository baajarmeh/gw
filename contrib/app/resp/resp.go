package resp

import (
	"bytes"
	"encoding/json"
	"io"
)

type IApiResp interface {
	GetCode() int
	GetContentType() string
	GetContentLength() int64
	GetHeaders() map[string]string
	GetBodyReader() io.Reader
}

type ApiRespBase struct {
	Code          int
	ContentType   string
	ContentLength int64
	Headers       map[string]string
	BodyObj       interface{}
}

type ApiRespRaw struct {
	ApiRespBase
}

func (obj ApiRespRaw) GetCode() int {
	return obj.Code
}
func (obj ApiRespRaw) GetContentType() string {
	return obj.ContentType
}
func (obj ApiRespRaw) GetContentLength() int64 {
	return obj.ContentLength
}
func (obj ApiRespRaw) GetHeaders() map[string]string {
	return obj.Headers
}
func (obj ApiRespRaw) GetBodyReader() io.Reader {
	panic("not implements.")
}

type ApiRespErrNotFound struct {
	ApiRespBase
}

type ApiRespErrForbidden struct {
	ApiRespBase
}

type ApiRespErrInternalError struct {
	ApiRespBase
}

type ApiRespJSON struct {
	ApiRespBase
}
func (obj ApiRespJSON) GetCode() int {
	return obj.Code
}
func (obj ApiRespJSON) GetContentType() string {
	return obj.ContentType
}
func (obj ApiRespJSON) GetContentLength() int64 {
	return obj.ContentLength
}
func (obj ApiRespJSON) GetHeaders() map[string]string {
	return obj.Headers
}
func (obj ApiRespJSON) GetBodyReader() io.Reader {
	body, _ := json.Marshal(obj.BodyObj)
	var r io.Reader
	var buf bytes.Buffer
	r = &buf
	buf.Write(body)
	return r
}

type ApiRespPlainText struct {
	ApiRespBase
}

func HttpApiRespRaw(code int, body interface{}) *ApiRespRaw {
	return &ApiRespRaw{}
}

func HttpApiRespJSON(code int, body interface{}) *ApiRespJSON {
	obj := &ApiRespJSON{}
	obj.Code = code
	obj.ContentType = "application/json"
	obj.BodyObj = body
}
