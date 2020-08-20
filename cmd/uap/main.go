package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap"
)

func main() {
	server := gw.DefaultServer()
	server.Patch(uap.New())
	//gw.Run(server)
	server.Serve()
}
