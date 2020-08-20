package generator

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/generator/services"
)

type App struct {
}

func (a App) Name() string {
	return "generator"
}

func (a App) Router() string {
	return "/gw/" + a.Name()
}

func (a App) Register(router *gw.RouterGroup) {
	router.GET("create-js", services.CreateJS)
}

func (a App) Migrate(ctx gw.MigrationContext) {
}

func (a App) Use(option *gw.ServerOption) {
}

func (a App) OnStart(state gw.ServerState) {

}

func (a App) OnShutDown(state gw.ServerState) {

}
