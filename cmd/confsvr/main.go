package main

import (
	"github.com/oceanho/gw/contrib/app"
	appConf "github.com/oceanho/gw/contrib/app/conf"
	"github.com/oceanho/gw/web/apps/confsvr"
)

func main() {
	bc := appConf.DefaultBootStrapConfig()
	conf := app.NewServerOption(bc)
	conf.Name = "confsvr"
	conf.Addr = ":8090"
	conf.Mode = "release"

	server := app.New(conf)
	registerApps(server)
	server.Serve()
}

func registerApps(server *app.ApiHostServer) {
	server.Register(confsvr.New())
}
