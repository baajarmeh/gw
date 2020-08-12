package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/contrib/apps/confsvr"
)

func main() {
	bsc := conf.DefaultBootConfig()
	opts := gw.NewServerOption(bsc)
	opts.Name = "confsvr"
	opts.Addr = ":8090"

	server := gw.NewServerWithOption(opts)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	server.Register(confsvr.New())
}
