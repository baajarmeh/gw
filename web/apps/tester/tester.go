package tester

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/tester/api"
)

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (u App) Name() string {
	return "gw.tester"
}

func (u App) BaseRouter() string {
	return "tester"
}

func (u App) Register(router *gw.RouteGroup) {
	router.GET("test/200", api.GetTester)
	router.GET("test/500", api.GetTester500)
}

func (u App) Migrate(store gw.Store) {
	// db := store.GetDbStore()
}
