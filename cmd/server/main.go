package main

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/stor"
	"github.com/oceanho/gw/web/apps/uap"
)

func main() {
	server := app.New()
	registerApps(server)
	server.Serve()
}

func registerApps(server *app.ApiServer) {
	server.Register(uap.New())
	server.Register(stor.New())
}
