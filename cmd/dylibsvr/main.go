package main

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/contrib/app/conf"
)

func main() {
	bcs := conf.LoadBootStrapConfigFromFile("config/boot.yaml")
	conf := app.NewServerOption(bcs)
	conf.Name = "confsvr"
	conf.Addr = ":8090"
	conf.Mode = "debug"

	server := app.New(conf)
	registerApps(server)
	server.Serve()
}

func registerApps(server *app.ApiHostServer) {
	dir := "./build/dylib"
	server.RegisterByPluginDir(dir)
}
