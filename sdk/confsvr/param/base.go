package param

import (
	json "github.com/json-iterator/go"
	"strings"
)

type IRequest interface {
	Method() string
	Url() string
	Body() []byte
	Headers() map[string]string
}

type Resp struct {
	Status string `json:"status"`
}

func (resp Resp) IsSucc() bool {
	return succFlags[strings.ToLower(resp.Status)]
}

func (resp Resp) String() string {
	b, _ := json.Marshal(resp)
	return string(b)
}

func IsHttpSucc(httpCode int) bool {
	return succCodes[httpCode]
}

var succFlags map[string]bool
var succCodes map[int]bool

func init() {
	succFlags = make(map[string]bool)
	succFlags["yes"] = true
	succFlags["ok"] = true
	succFlags["succ"] = true
	succFlags["true"] = true
	succFlags["y"] = true
	succCodes = make(map[int]bool)
	succCodes[200] = true
	succCodes[201] = true
	succCodes[202] = true
	succCodes[203] = true
}
