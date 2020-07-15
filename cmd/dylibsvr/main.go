package main

import (
	"github.com/oceanho/gw/contrib/app"
)

func main() {
	conf := app.NewServerOption()
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
