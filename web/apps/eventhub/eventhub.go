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

func (u App) Name() string {
	return "gw.event-hub"
}

func (u App) Router() string {
	return "event-hub"
}

func (u App) Register(router *gw.RouterGroup) {

}

func (u App) Migrate(store gw.Store) {
	// db := store.GetDbStore()
}
