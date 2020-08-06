package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
)

func main() {
	bcs := conf.DefaultBootConfig()
	opts := gw.NewServerOption(bcs)
	server := gw.New(opts)
	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	// register Your app at here.
	//server.Register(myapp.App{})
}
