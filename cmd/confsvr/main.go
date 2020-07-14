package main

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/confsvr"
)

func main() {
	conf := app.NewOption()
	conf.Name = "confsvr"
	conf.Addr = ":8090"
	conf.Mode = "debug"

	server := app.New(conf)
	registerApps(server)
	server.Serve()
}

func registerApps(server *app.ApiServer) {
	server.Register(confsvr.New())
}
