package confsvr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/oceanho/gw/sdk/confsvr/param"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	options *Option
}

type Option struct {
	Addr               string
	Service            string
	Version            string
	Proto              string
	MaxRetries         int
	QueryStateInterval int
}

func (c *Client) Do(req param.IRequest, output interface{}) (int, error) {
	var reader io.Reader
	var buffer bufio.ReadWriter
	body := req.Body()
	if len(body) > 0 {
		reader = &buffer
		buffer.Write(req.Body())
	}

	uri := req.Url()
	uri = strings.TrimLeft(uri, "/")
	uri = strings.TrimRight(uri, "/")

	url := fmt.Sprintf("%s/%s/%s/%s", c.options.Addr, c.options.Version, c.options.Service, uri)
	oriReq, _ := http.NewRequest(req.Method(), url, reader)
	for k, header := range req.Headers() {
		oriReq.Header.Set(k, header)
	}

	client := http.DefaultClient
	resp, err := client.Do(oriReq)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("do http request: %v", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read resp body: %v", err)
	}
	if !param.IsHttpSucc(resp.StatusCode) {
		return resp.StatusCode, fmt.Errorf("status code [%d] are not excepted, resp text: \"%s\"", resp.StatusCode, string(b))
	}
	err = json.Unmarshal(b, output)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("decode str to json: \"%s\", %v", string(b), err)
	}
	return resp.StatusCode, nil
}

var (
	defaultAddr         = "https://127.0.0.1:8080"
	defaultVersion      = "//v1"
	defaultService      = "confsvr"
	defaultProto        = "http"
	defaultMaxRetries   = 10
	defaultPullInterval = 15
)

func DefaultOptions() *Option {
	return &Option{
		Addr:               defaultAddr,
		Proto:              defaultProto,
		Service:            defaultService,
		Version:            defaultVersion,
		MaxRetries:         defaultMaxRetries,
		QueryStateInterval: defaultPullInterval,
	}
}

func NewClient(opts *Option) *Client {
	opts.Addr = strings.TrimRight(opts.Addr, "/")
	opts.Version = strings.TrimLeft(opts.Version, "/")
	opts.Version = strings.TrimRight(opts.Version, "/")
	client := &Client{
		options: opts,
	}
	return client
}
