package main

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap"
)

func main() {
	server := gw.DefaultServer()
	server.Register(uap.New())
	server.Serve()
}
