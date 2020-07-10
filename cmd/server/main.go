package main

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/stor"
	"github.com/oceanho/gw/web/apps/uap"
)

func main() {
	server := app.New()
	server.Register(uap.New())
	server.Register(stor.New())
	server.Serve()
}
