package admin

import "github.com/oceanho/gw"

type App struct {
}

func (a App) Name() string {
	return "gw.admin"
}

func (a App) Router() string {
	return "gw/admin"
}

func (a App) Register(router *gw.RouterGroup) {
	//gw.RegisterRestAPIs(router,)
}

func (a App) Migrate(ctx gw.MigrationContext) {
	// Nothing to do.
}

func (a App) Use(opt *gw.ServerOption) {
}
