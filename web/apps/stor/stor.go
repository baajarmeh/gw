package stor

import (
	gw2 "github.com/oceanho/gw"
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
	return "oceanho.stor"
}

func (u App) BaseRouter() string {
	return "stor"
}

func (u App) Register(router *gw2.ApiRouteGroup) {
	router.GET("object/create", api.CreateObject)
	router.POST("object/modify", api.ModifyObject)
}
