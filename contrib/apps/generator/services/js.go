package services

import "github.com/oceanho/gw"

func CreateJS(ctx *gw.Context) {
	//s := ctx.GetHostServer()
	//routers := s.GetRouters()
	js := "<javascript>alert('hello')</javascript>"
	s := gw.GetHostServer(ctx)
	ctx.Context.String(200, js)
	ctx.Context.String(200, s.Name)
}
