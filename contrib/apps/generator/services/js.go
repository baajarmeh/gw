package services

import "github.com/oceanho/gw"

func CreateJS(ctx *gw.Context)  {
	//s := ctx.GetHostServer()
	//routers := s.GetRouters()
	js := string("<javascript>alert('hello')</javascript>")
	ctx.Context.String(200,js)
}