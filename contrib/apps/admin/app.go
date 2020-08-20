package admin

import "github.com/oceanho/gw"

type App struct {
}

func (a App) Name() string {
	return "gw.admin"
}

func (a App) Router() string {
	return "gw/admin"
}

func (a App) Register(router *gw.RouterGroup) {
}

func (a App) Use(opt *gw.ServerOption) {
}

func (a App) Migrate(state gw.ServerState) {
	// Nothing to do.
}

func (a App) OnStart(state gw.ServerState) {

}

func (a App) OnShutDown(state gw.ServerState) {

}
