package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/web/apps/confsvr"
)

func main() {
	bc := conf.DefaultBootStrapConfig()
	conf := gw.NewServerOption(bc)
	conf.Name = "confsvr"
	conf.Addr = ":8090"
	conf.Mode = "release"

	server := gw.New(conf)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	server.Register(confsvr.New())
}
