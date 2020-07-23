package stor

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/stor/api"
)

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (u App) Name() string {
	return "gw.event-hub"
}

func (u App) BaseRouter() string {
	return "event-hub"
}

func (u App) Register(router *gw.RouteGroup) {
	router.GET("object/create", api.CreateObject)
	router.POST("object/modify", api.ModifyObject)
}
