package tester

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/tester/api"
	"github.com/oceanho/gw/web/apps/tester/dto"
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

	router.GET("test/create", api.CreateMyTester)
	router.GET("test/query", api.QueryMyTester)

	router.GET("test/200", api.GetTester)

	router.GET("test/400", api.GetTester400)
	router.GET("test/400-err", api.GetTester400WithCustomErr)
	router.GET("test/400-payload", api.GetTester400WithCustomPayload)
	router.GET("test/400-err-payload", api.GetTester400WithCustomPayloadErr)

	router.GET("test/401", api.GetTester401)
	router.GET("test/401-err", api.GetTester401WithCustomErr)
	router.GET("test/401-payload", api.GetTester401WithCustomPayload)
	router.GET("test/401-err-payload", api.GetTester401WithCustomPayloadErr)

	router.GET("test/403", api.GetTester403)
	router.GET("test/403-err", api.GetTester403WithCustomErr)
	router.GET("test/403-payload", api.GetTester403WithCustomPayload)
	router.GET("test/403-err-payload", api.GetTester403WithCustomPayloadErr)

	router.GET("test/404", api.GetTester404)
	router.GET("test/404-err", api.GetTester404WithCustomErr)
	router.GET("test/404-payload", api.GetTester404WithCustomPayload)
	router.GET("test/404-err-payload", api.GetTester404WithCustomPayloadErr)

	router.GET("test/500", api.GetTester500)
	router.GET("test/500-err", api.GetTester500WithCustomErr)
	router.GET("test/500-payload", api.GetTester500WithCustomPayload)
	router.GET("test/500-err-payload", api.GetTester500WithCustomPayloadErr)
}

func (u App) Migrate(store gw.Store) {
	db := store.GetDbStore()
	db.AutoMigrate(&dto.MyTester{})
}
