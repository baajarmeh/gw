package confsvr

import (
	"github.com/oceanho/gw/sdk/confsvr/req"
	"testing"
)

func TestClient_Do(t *testing.T) {
	opts := DefaultOptions()
	opts.Addr = "http://127.0.0.1:8090/"
	client := NewClient(opts)
	reqobj := req.GetAuthRequest{
		AccessKeyId:     "NX3JM7ERkbKEJ3kujbvVeqdUhavpFFPf",
		AccessKeySecret: "7t3JAuwwfTgUfjkbJVoY3KJg3uLPrWv4Kpvjnv9MRLtTcJcTsHxJTg9THMEFpe7U",
	}
	respobj := &req.RespGetAuth{}
	code, err := client.Do(reqobj, respobj)
	if err != nil {
		t.Fatalf("client.Do(reqobj, respobj), status code: %d, resp: %v", code, err)
	}
	t.Logf("result is: %v", respobj)
}
