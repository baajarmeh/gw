package api

import (
	"github.com/oceanho/gw"
)

func GetTester(c *gw.Context) {
	c.OK(struct {
		RequestID string
	}{
		RequestID: c.RequestID,
	})
}

func GetTester500(c *gw.Context) {
	c.Err500(4000)
}
