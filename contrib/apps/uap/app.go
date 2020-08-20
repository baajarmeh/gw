package uap

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"github.com/oceanho/gw/contrib/apps/uap/gwImpls"
	"github.com/oceanho/gw/contrib/apps/uap/restAPIs"
)

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (a App) Name() string {
	return "gw.uap"
}

func (a App) Router() string {
	return "uap"
}

func (a App) Register(router *gw.RouterGroup) {
	router.RegisterRestAPIs(restAPIs.DynamicAPIs()...)
}

func (a App) Migrate(ctx gw.MigrationContext) {
	db := ctx.Store.GetDbStore()
	_ = db.AutoMigrate(&dbModel.User{}, &dbModel.Role{}, &dbModel.UserProfile{}, &dbModel.UserRoleMapping{})
}

func (a App) Use(opt *gw.ServerOption) {
	opt.AuthManagerHandler = func(state gw.ServerState) gw.IAuthManager {
		return gwImpls.DefaultAuthManager(state)
	}
	opt.PermissionManagerHandler = func(state gw.ServerState) gw.IPermissionManager {
		return gwImpls.DefaultPermissionManager(state)
	}
}
