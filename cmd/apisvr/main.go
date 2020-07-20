package main

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/stor"
	"github.com/oceanho/gw/web/apps/uap"
)

func main() {
	server := app.Default()
	registerApps(server)
	server.Serve()
}

func registerApps(server *app.ApiHostServer) {
	server.Register(uap.New())
	server.Register(stor.New())
}