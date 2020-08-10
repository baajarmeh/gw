package stor

import (
	"github.com/oceanho/gw"
)

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (a App) Name() string {
	return "gw.event-hub"
}

func (a App) Router() string {
	return "event-hub"
}

func (a App) Register(router *gw.RouterGroup) {
}

func (a App) Migrate(ctx gw.MigrationContext) {
	// db := store.GetDbStore()
}

func (a App) Use(opt *gw.ServerOption) {
}
