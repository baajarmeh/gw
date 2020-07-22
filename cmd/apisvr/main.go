package main

import (
	"github.com/oceanho/gw/contrib/app"
	appConf "github.com/oceanho/gw/contrib/app/conf"
	"github.com/oceanho/gw/web/apps/stor"
	"github.com/oceanho/gw/web/apps/uap"
)

func main() {
	// server := app.Default()

	bc := appConf.DefaultBootStrapConfig()
	conf := app.NewServerOption(bc)
	conf.Name = "confsvr"
	conf.Addr = ":8080"
	conf.Mode = "release"
	server := app.New(conf)
	registerApps(server)
	server.Serve()
}

func registerApps(server *app.ApiHostServer) {
	server.Register(uap.New())
	server.Register(stor.New())
}
