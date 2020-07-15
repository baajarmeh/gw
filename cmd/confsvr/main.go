package main

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/confsvr"
)

func main() {
	conf := app.NewServerOption()
	conf.Name = "confsvr"
	conf.Addr = ":8090"
	conf.Mode = "debug"

	server := app.New(conf)
	registerApps(server)
	server.Serve()
}

func registerApps(server *app.ApiHostServer) {
	server.Register(confsvr.New())
}
