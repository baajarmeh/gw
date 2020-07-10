package uap

import "github.com/gin-gonic/gin"

func init() {
}

type App struct {
}

func New() *App {
	return &App{}
}

func (u App) Name() string {
	return "oceanho.uap"
}

func (u App) Register(router *gin.Engine) {

}
