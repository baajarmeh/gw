package api

import "github.com/oceanho/gw"

func Err401(c *gw.Context) {
	c.JSON401(401)
}

func Err500(c *gw.Context) {
	panic("err 500.")
}
