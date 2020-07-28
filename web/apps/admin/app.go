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
	//gw.RegisterControllers(router,)
}

func (a App) Migrate(store gw.Store) {
	// Nothing to do.
}
