package main

import (
	"github.com/oceanho/gw/contrib/apps/uap"
)

var AppPlugin uap.App

func init() {
	AppPlugin = uap.New()
}
