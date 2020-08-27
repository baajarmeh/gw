package gw

import (
	"github.com/gin-gonic/gin"
)

// App represents a application.
//
// The APIs called by gw framework order as
//
// 1. Use()
//
// 2. Register()
//
// 3. Migrate()
//
// 4. OnStart()
//
// 5. OnShutDown
//
type App interface {

	// Name define a API that return as your app name.
	Name() string

	// Router define a API that should be return your app base route path.
	// It's will be used to create a new *RouteGroup object and that used by Register(...) .
	Router() string

	// Register define a API that for register your app router inside.
	Register(router *RouterGroup)

	// Use define a API that for modify Server Options capability AT your Application.
	Use(option *ServerOption)

	// Migrate define a API that for create your app's database migrations and permission initialization inside.
	Migrate(state *ServerState)

	// OnStart define a API that notify your Application when server starting before.
	OnStart(state *ServerState)

	// OnShutDown define a API that notify your Application when server shutdown before.
	OnShutDown(state *ServerState)

	// Fixme(OceanHo): need?
	// Use define a API that tells gw framework, Your application Depend On apps.
	//DependOn(registeredApps map[string]App) []string
}

type internalApp struct {
	instance    App
	isPatchOnly bool
}

// MigrationContext represents a Migration Context Object.
type MigrationContext struct {
	ServerState
}

// IDynamicRestAPI represents a Dynamic Rest Style API interface.
type IDynamicRestAPI interface {
	// Name define a API that returns Your RestAPI name(such as resource name)
	// It's will be used as router prefix.
	Name() string
}

// ErrorHandler represents a http Error handler.
type ErrorHandler func(requestId string, httpRequest string, headers []string, stack string, err []*gin.Error)
