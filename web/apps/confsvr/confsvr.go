package confsvr

import (
	gw2 "github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/confsvr/"
)

type App struct {
}

func New() App {
	return App{}
}

func (u App) Name() string {
	return "oceanho.confsvr"
}

func (u App) BaseRouter() string {
	return "confsvr"
}

func (u App) Register(router *gw2.RouteGroup) {
	// Auth service routers.
	router.GET("auth/auth", api.GetAuth)
	router.GET("auth/create", api.CreateEnv)
	router.GET("auth/modify", api.ModifyAuth)
	router.GET("auth/destroy", api.DestroyAuth)

	// Env service routers.
	router.GET("env/get", api.GetEnv)
	router.GET("env/create", api.CreateEnv)
	router.GET("env/modify", api.ModifyEnv)
	router.GET("env/destroy", api.DestroyEnv)

	// NameSpace service routers.
	router.GET("ns/get", api.GetNS)
	router.GET("ns/create", api.CreateNS)
	router.GET("ns/modify", api.ModifyNS)
	router.GET("ns/destroy", api.DestroyNS)
}
