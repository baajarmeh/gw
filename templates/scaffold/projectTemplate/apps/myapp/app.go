package myapp

import (
	"github.com/oceanho/gw"
)

//
// Impl gw.App interface at here.
// reference: https://github.com/oceanho/gw/wiki/Scaffold-Guides#6-appgo
//

type App struct {
}

func New() App {
	return App{}
}

func (a App) Name() string {
	// return your app name
	// Example: return "gw.uap"
	return "my-app"
}

func (a App) Router() string {
	// return your app router group name.
	return "my-app"
}

func (a App) Register(router *gw.RouterGroup) {
	// your app routers
}

func (a App) Migrate(state *gw.ServerState) {
	// your database migrations
}

func (a App) Use(opt *gw.ServerOption) {
	// modify serverOptions if your have necessary.
}
