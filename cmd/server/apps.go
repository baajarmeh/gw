package main

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/stor"
	"github.com/oceanho/gw/web/apps/uap"
)

func registerApps(server *app.ApiServer) {
	server.Register(uap.New())
	server.Register(stor.New())
}
