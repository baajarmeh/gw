package pvm

import "github.com/oceanho/gw"

type App struct {
	name            string
	router          string
	registerFunc    func(router *gw.RouterGroup)
	useFunc         func(option *gw.ServerOption)
	migrateFunc     func(state *gw.ServerState)
	onStartFunc     func(state *gw.ServerState)
	onShoutDownFunc func(state *gw.ServerState)
}

func New() gw.App {
	var app = &App{
		name:   "pvm",
		router: "oceanho/gw-pvm",
		registerFunc: func(router *gw.RouterGroup) {

		},
		useFunc: func(option *gw.ServerOption) {

		},
		migrateFunc: func(state *gw.ServerState) {

		},
		onStartFunc: func(state *gw.ServerState) {

		},
		onShoutDownFunc: func(state *gw.ServerState) {

		},
	}
	return app
}

func (a App) Name() string {
	return a.name
}

func (a App) Router() string {
	return a.router
}

func (a App) Register(router *gw.RouterGroup) {
	a.registerFunc(router)
}

func (a App) Use(option *gw.ServerOption) {
	a.useFunc(option)
}

func (a App) Migrate(state *gw.ServerState) {
	a.migrateFunc(state)
}

func (a App) OnStart(state *gw.ServerState) {
	a.onStartFunc(state)
}

func (a App) OnShutDown(state *gw.ServerState) {
	a.onShoutDownFunc(state)
}
