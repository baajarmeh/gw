package main

import (
	"github.com/oceanho/gw/contrib/gwcli"
	"os"
)


func main() {
	app := gwcli.App()
	app.Run(os.Args)
}
