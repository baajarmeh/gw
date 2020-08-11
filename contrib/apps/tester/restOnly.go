package tester

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/tester/rest"
)

func init() {
}

type AppRestOnly struct {
}

func NewAppRestOnly() AppRestOnly {
	return AppRestOnly{}
}

func (a AppRestOnly) Name() string {
	return "gw-rest-app-only"
}

func (a AppRestOnly) Router() string {
	return "router-only-tester"
}

func (a AppRestOnly) Register(router *gw.RouterGroup) {
	gw.RegisterRestAPI(router, &rest.MyTesterRestAPI{})
}

func (a AppRestOnly) Migrate(ctx gw.MigrationContext) {
	//db := ctx.Store.GetDbStore()
	//pm := ctx.PermissionManager
	perms := gw.NewPermAll("MyTester")
	err := ctx.PermissionManager.Create("mytester", perms...)
	if err != nil {
		panic(err)
	}
}

func (a AppRestOnly) Use(opt *gw.ServerOption) {
}
