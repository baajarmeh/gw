package uap

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/uap/api"
)

func init() {
}

type App struct {
}

func New() *App {
	return &App{}
}

func (u App) Name() string {
	return "oceanho.uap"
}

func (u App) BaseRouter() string {
	return "uap"
}

func (u App) Register(router *app.ApiRouteGroup) {
	router.GET("tenant/query", api.GetTenant)
}
