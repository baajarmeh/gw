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
	Register(router *gin.RouterGroup)
}

type Router struct {
	gin.Engine
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
	httpServer            *gin.Engine
	router                *gin.RouterGroup
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
	apiServer := &ApiServer{
		Addr:                  appDefaultAddr,
		Name:                  appDefaultName,
		Restart:               appDefaultRestart,
		Mode:                  appDefaultMode,
		StartBeforeHandler:    appDefaultStartBeforeHandler,
		ShutDownBeforeHandler: appDefaultShutDownBeforeHandler,
		httpServer:            engine,
		ApiPrefix:             appDefaultApiPrefix,
		ApiVersion:            appDefaultApiVersion,
		apps:                  make(map[string]App),
	}
	apiGroupPath := fmt.Sprintf("%s/%s", apiServer.ApiPrefix, apiServer.ApiVersion)
	apiServer.router = apiServer.httpServer.Group(apiGroupPath)
	return apiServer
}

func (server *ApiServer) Register(apps ...App) {
	server.locker.Lock()
	defer server.locker.Unlock()
	for _, app := range apps {
		appName := app.Name()
		if _, ok := server.apps[appName]; !ok {
			server.apps[appName] = app
			rg := server.router.Group(app.BaseRouter())
			app.Register(rg)
		}
	}
}

func (server *ApiServer) Serve() {
	sigs := make(chan os.Signal, 1)
	handler := server.StartBeforeHandler
	if handler != nil {
		err := server.StartBeforeHandler(server)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	err := server.httpServer.Run(server.Addr)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	handler = server.ShutDownBeforeHandler
	if handler != nil {
		err := server.ShutDownBeforeHandler(server)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
}
