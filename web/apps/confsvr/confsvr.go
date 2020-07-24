package confsvr

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/confsvr/api"
)

type App struct {
}

func New() App {
	return App{}
}

func (u App) Name() string {
	return "gw.confsvr"
}

func (u App) Router() string {
	return "confsvr"
}

func (u App) Register(router *gw.RouteGroup) {
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

func (u App) Migrate(store gw.Store) {
	// db := store.GetDbStore()
}
