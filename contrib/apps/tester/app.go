package tester

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/tester/api"
	"github.com/oceanho/gw/contrib/apps/tester/dto"
	"github.com/oceanho/gw/contrib/apps/tester/infra"
	"github.com/oceanho/gw/contrib/apps/tester/rest"
	"gorm.io/gorm"
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

	router.POST("test/create", api.CreateMyTester, infra.DecoratorList.Creation())
	router.GET("test/query", api.QueryMyTester, gw.NewDecorators(gw.NewStoreDbSetupDecorator(func(ctx gw.Context, db *gorm.DB) *gorm.DB {
		return db.Where("id >= 10")
	})).Append(gw.NewPermAllDecorator("a").All()...).All()...)

	router.GET("test/200", api.GetTester)

	router.GET("test/400", api.GetTester400)
	router.GET("test/400-err", api.GetTester400WithCustomErr)
	router.GET("test/400-payload", api.GetTester400WithCustomPayload)
	router.GET("test/400-err-payload", api.GetTester400WithCustomPayloadErr)

	router.GET("test/401", api.GetTester401)
	router.GET("test/401-err", api.GetTester401WithCustomErr)
	router.GET("test/401-payload", api.GetTester401WithCustomPayload)
	router.GET("test/401-err-payload", api.GetTester401WithCustomPayloadErr)

	router.GET("test/403", api.GetTester403, gw.NewPermAllDecorator("test-403").All()...)
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

	router.RegisterRestAPIs(&rest.MyTesterRestAPI{})

	router.GET("err/401", api.Err401)
	router.GET("err/500", api.Err500)
}

func (a App) Migrate(state gw.ServerState) {
	db := state.Store().GetDbStore()
	db.AutoMigrate(&dto.MyTester{})
}

func (a App) Use(opt *gw.ServerOption) {

}

func (a App) OnStart(state gw.ServerState) {

}

func (a App) OnShutDown(state gw.ServerState) {

}
