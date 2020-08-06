package tester

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/tester/api"
	"github.com/oceanho/gw/web/apps/tester/dto"
	"github.com/oceanho/gw/web/apps/tester/infra"
	"github.com/oceanho/gw/web/apps/tester/rest"
)

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (a App) Name() string {
	return "gw.tester"
}

func (a App) Router() string {
	return "tester"
}

func (a App) Register(router *gw.RouterGroup) {

	router.POST("test/create", api.CreateMyTester, infra.DecoratorList.ReadAll())
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

	gw.RegisterRestAPI(router, &rest.MyTesterRestAPI{})

	router.GET("err/401", api.Err401)
	router.GET("err/500", api.Err500)
}

func (a App) Migrate(store gw.Store) {
	db := store.GetDbStore()
	db.AutoMigrate(&dto.MyTester{})
}

func (a App) Use(opt *gw.ServerOption) {
}
