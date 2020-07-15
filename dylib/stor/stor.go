package stor

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/stor/api"
)

var AppPlugin App

func init() {
	AppPlugin = New()
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

func (u App) Register(router *app.ApiRouteGroup) {
	router.GET("/get", api.CreateObject)
	router.POST("/modify", api.ModifyObject)
}
