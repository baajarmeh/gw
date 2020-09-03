package gw

import (
	"github.com/gin-gonic/gin"
)

// AppInfo ...
type AppInfo struct {
	Name       string
	Router     string
	Descriptor string
}

// IAppManager ...
type IAppManager interface {
	Create(app AppInfo) error
	QueryByName(name string) *AppInfo
}

func EmptyAppManager(state *ServerState) IAppManager {
	return EmptyAppManagerImpl{
		state: state,
	}
}

type EmptyAppManagerImpl struct {
	state *ServerState
}

func (d EmptyAppManagerImpl) Create(app AppInfo) error {
	return nil
}

func (d EmptyAppManagerImpl) QueryByName(name string) *AppInfo {
	return nil
}

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

	// Register define a API that it returns your app meta info.
	Meta() AppInfo

	// Register define a API that for register your app router inside.
	Register(router *RouterGroup)

	// Use define a API that for modify server options capability AT your Application.
	Use(option *ServerOption)

	// OnPrepare define a API that for create your app's db migration,some initializations.
	OnPrepare(state *ServerState)

	// OnStart define a API that notify your Application when server starting before.
	OnStart(state *ServerState)

	// OnShutDown define a API that notify your Application when server shutdown before.
	OnShutDown(state *ServerState)
}

type internalApp struct {
	instance    App
	isPatchOnly bool
}

// IDynamicRestAPI represents a Dynamic Rest Style API interface.
type IDynamicRestAPI interface {
	// Name define a API that returns Your RestAPI name(such as resource name)
	// It's will be used as router prefix.
	Name() string
}

// ErrorHandler represents a http Error handler.
type ErrorHandler func(requestId string, httpRequest string, headers []string, stack string, err []*gin.Error)
