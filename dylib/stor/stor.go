package main

import (
	"github.com/oceanho/gw/contrib/apps/stor"
)

var AppPlugin stor.App

func init() {
	AppPlugin = stor.New()
}
