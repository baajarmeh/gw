package main

import (
	"github.com/oceanho/gw/web/apps/stor"
)

var AppPlugin stor.App

func init() {
	AppPlugin = stor.New()
}
