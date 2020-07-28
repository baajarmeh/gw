package confsvr

import (
	"github.com/oceanho/gw/sdk/confsvr/param"
	"testing"
)

func TestClient_Do(t *testing.T) {
	opts := DefaultOptions()
	opts.Addr = "http://127.0.0.1:8090/"
	client := NewClient(opts)
	req := param.ReqGetAuth{
		AccessKeyId:     "NX3JM7ERkbKEJ3kujbvVeqdUhavpFFPf",
		AccessKeySecret: "7t3JAuwwfTgUfjkbJVoY3KJg3uLPrWv4Kpvjnv9MRLtTcJcTsHxJTg9THMEFpe7U",
	}
	resp := &param.RspGetAuth{}
	code, err := client.Do(req, resp)
	if err != nil {
		t.Fatalf("client.Do(req, respobj), status code: %d, resp: %v", code, err)
	}
	t.Logf("result is: %v", resp)
}
