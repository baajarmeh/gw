package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/contrib/apps/generator"
)

func main() {
	bcs := conf.NewBootConfigFromFile("config/apps/generator/boot.yaml")
	opts := gw.NewServerOption(bcs)
	opts.Name = "my-generator-api"
	server := gw.NewServerWithOption(opts)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	server.Register(generator.App{})
}
