package main

import (
	"github.com/oceanho/gw/contrib/app"
)

func main() {
	server := app.New()
	registerApps(server)
	server.Serve()
}
