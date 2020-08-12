package main

import (
	"github.com/oceanho/gw"
	conf2 "github.com/oceanho/gw/conf"
)

func main() {
	bcs := conf2.NewBootConfigFromFile("config/boot.yaml")
	opts := gw.NewServerOption(bcs)
	opts.Name = "confsvr"
	opts.Addr = ":8090"

	server := gw.NewServerWithOption(opts)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	dir := "./build/dylib"
	server.RegisterByPluginDir(dir)
}
