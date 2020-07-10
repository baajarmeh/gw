package stor

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/web/apps/stor/api"
)

func init() {
}

type App struct {
}

func New() *App {
	return &App{}
}

func (u App) Name() string {
	return "oceanho.stor"
}

func (u App) BaseRouter() string {
	return "stor"
}

func (u App) Register(router *gin.RouterGroup) {
	router.GET("obj/<uid:int64>", api.CreateObject)
}
