package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/contrib/apps/stor"
	"github.com/oceanho/gw/contrib/apps/uap"
)

func main() {
	bcs := conf.DefaultBootConfig()
	opts := gw.NewServerOption(bcs)
	opts.Name = "my api"
	server := gw.New(opts)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	server.Register(uap.New())
	server.Register(stor.New())
}
