package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App interface {
	Name() string
	BaseRouter() string
	Register(router *ApiRouteGroup)
}

type ApiServer struct {
	Addr                  string
	Name                  string
	Mode                  string
	Restart               string
	ApiPrefix             string
	ApiVersion            string
	StartBeforeHandler    func(server *ApiServer) error
	ShutDownBeforeHandler func(server *ApiServer) error
	locker                sync.Mutex
	router                *ApiRouter
	apps                  map[string]App
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "oceanho.app"
	appDefaultRestart               = "always"
	appDefaultMode                  = "release"
	appDefaultApiPrefix             = "api"
	appDefaultApiVersion            = "v1"
	appDefaultStartBeforeHandler    = func(server *ApiServer) error { return nil }
	appDefaultShutDownBeforeHandler = func(server *ApiServer) error { return nil }
)

func New() *ApiServer {
	engine := gin.New()
	httpRouter := &ApiRouter{
		server: engine,
	}
	apiServer := &ApiServer{
		Addr:                  appDefaultAddr,
		Name:                  appDefaultName,
		Restart:               appDefaultRestart,
		Mode:                  appDefaultMode,
		StartBeforeHandler:    appDefaultStartBeforeHandler,
		ShutDownBeforeHandler: appDefaultShutDownBeforeHandler,
		ApiPrefix:             appDefaultApiPrefix,
		ApiVersion:            appDefaultApiVersion,
		apps:                  make(map[string]App),
	}
	httpRouter.router = httpRouter.server.Group(fmt.Sprintf("%s/%s", apiServer.ApiPrefix, apiServer.ApiVersion))
	apiServer.router = httpRouter
	return apiServer
}

func (apiServer *ApiServer) Register(apps ...App) {
	apiServer.locker.Lock()
	defer apiServer.locker.Unlock()
	for _, app := range apps {
		appName := app.Name()
		if _, ok := apiServer.apps[appName]; !ok {
			apiServer.apps[appName] = app
			apiRg := apiServer.router.Group(app.BaseRouter())
			app.Register(apiRg)
		}
	}
}

func (apiServer *ApiServer) Serve() {
	sigs := make(chan os.Signal, 1)
	handler := apiServer.StartBeforeHandler
	if handler != nil {
		err := apiServer.StartBeforeHandler(apiServer)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	err := apiServer.router.server.Run(apiServer.Addr)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	handler = apiServer.ShutDownBeforeHandler
	if handler != nil {
		err := apiServer.ShutDownBeforeHandler(apiServer)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
}
