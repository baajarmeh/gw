package main

import (
	"github.com/oceanho/gw/web/apps/uap"
)

var AppPlugin uap.App

func init() {
	AppPlugin = uap.New()
}
