package main

import (
	"github.com/oceanho/gw"
	conf2 "github.com/oceanho/gw/conf"
)

func main() {
	bcs := conf2.LoadBootConfigFromFile("config/boot.yaml")
	conf := gw.NewServerOption(bcs)
	conf.Name = "confsvr"
	conf.Addr = ":8090"
	conf.Mode = "debug"

	server := gw.New(conf)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	dir := "./build/dylib"
	server.RegisterByPluginDir(dir)
}
