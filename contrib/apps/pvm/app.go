package pvm

import "github.com/oceanho/gw"

type App struct {
	name            string
	router          string
	onPrepareFunc   func(state *gw.ServerState)
	onStartFunc     func(state *gw.ServerState)
	onShoutDownFunc func(state *gw.ServerState)
	registerFunc    func(router *gw.RouterGroup)
	useFunc         func(option *gw.ServerOption)
}

const (
	appName   = "gw.pvm"
	appRouter = "oceanho/gw-pvm"
)

func New() gw.App {
	var app = &App{
		name:   appName,
		router: appRouter,
		registerFunc: func(router *gw.RouterGroup) {

		},
		useFunc: func(option *gw.ServerOption) {

		},
		onPrepareFunc: func(state *gw.ServerState) {

		},
		onStartFunc: func(state *gw.ServerState) {

		},
		onShoutDownFunc: func(state *gw.ServerState) {

		},
	}
	return app
}

func (a App) Meta() gw.AppInfo {
	return gw.AppInfo{
		Name:       a.name,
		Router:     a.router,
		Descriptor: "Gw Product Version Manager System.",
	}
}

func (a App) Register(router *gw.RouterGroup) {
	a.registerFunc(router)
}

func (a App) Use(option *gw.ServerOption) {
	a.useFunc(option)
}

func (a App) OnPrepare(state *gw.ServerState) {
	a.onPrepareFunc(state)
}

func (a App) OnStart(state *gw.ServerState) {
	a.onStartFunc(state)
}

func (a App) OnShutDown(state *gw.ServerState) {
	a.onShoutDownFunc(state)
}
