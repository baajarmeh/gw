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

func (a App) Register(router *gw.RouteGroup) {
	panic("implement me")
}

func (a App) Migrate(store gw.Store) {
	panic("implement me")
}
