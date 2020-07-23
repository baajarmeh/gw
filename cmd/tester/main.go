package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/web/apps/tester"
)

func main() {
	bcs := conf.DefaultBootStrapConfig()
	opts := gw.NewServerOption(bcs)
	opts.Name = "my-tester-api"
	server := gw.New(opts)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	server.Register(tester.New())
}
