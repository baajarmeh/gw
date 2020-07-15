package main

import (
	"github.com/oceanho/gw/web/apps/confsvr"
)

var AppPlugin confsvr.App

func init() {
	AppPlugin =  confsvr.New()
}
