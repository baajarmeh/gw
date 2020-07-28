package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/web/apps/confsvr"
)

func main() {
	bsc := conf.DefaultBootStrapConfig()
	conf := gw.NewServerOption(bsc)
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
