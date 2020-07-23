package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/web/apps/stor"
	"github.com/oceanho/gw/web/apps/uap"
)

func main() {
	// server := app.Default()

	bcs := conf.DefaultBootStrapConfig()
	opts := gw.NewServerOption(bcs)
	opts.Name = "my api"
	//opts.Mode = "release"
	server := gw.New(opts)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	server.Register(uap.New())
	server.Register(stor.New())
}
