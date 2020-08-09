package main

import (
	"github.com/oceanho/gw/contrib/apps/confsvr"
)

var AppPlugin confsvr.App

func init() {
	AppPlugin = confsvr.New()
}
