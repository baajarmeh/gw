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

type Config struct {
	Addr                  string
	Name                  string
	Mode                  string
	Restart               string
	ApiPrefix             string
	ApiVersion            string
	StartBeforeHandler    func(server *ApiServer) error
	ShutDownBeforeHandler func(server *ApiServer) error
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
	rootRouter			  *ApiRouter
	apps                  map[string]App
	conf                  *Config
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "oceanho.app"
	appDefaultRestart               = "always"
	appDefaultMode                  = "debug"
	appDefaultApiPrefix             = "api/v1"
	appDefaultApiVersion            = "Version 1.0"
	appDefaultStartBeforeHandler    = func(server *ApiServer) error { return nil }
	appDefaultShutDownBeforeHandler = func(server *ApiServer) error { return nil }
)

func NewConfig() *Config {
	conf := &Config{
		Addr:                  appDefaultAddr,
		Name:                  appDefaultName,
		Restart:               appDefaultRestart,
		Mode:                  appDefaultMode,
		StartBeforeHandler:    appDefaultStartBeforeHandler,
		ShutDownBeforeHandler: appDefaultShutDownBeforeHandler,
		ApiPrefix:             appDefaultApiPrefix,
		ApiVersion:            appDefaultApiVersion,
	}
	return conf
}

func Default() *ApiServer {
	return New(NewConfig())
}

func New(conf *Config) *ApiServer {
	gin.SetMode(conf.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	httpRouter := &ApiRouter{
		server: engine,
	}
	apiServer := &ApiServer{
		Addr:                  conf.Addr,
		Name:                  conf.Name,
		Restart:               conf.Restart,
		Mode:                  conf.Mode,
		StartBeforeHandler:    conf.StartBeforeHandler,
		ShutDownBeforeHandler: conf.ShutDownBeforeHandler,
		ApiPrefix:             conf.ApiPrefix,
		ApiVersion:            conf.ApiVersion,
		apps:                  make(map[string]App),
		conf:                  conf,
	}
	httpRouter.router = httpRouter.server.Group(apiServer.ApiPrefix)
	apiServer.rootRouter = &ApiRouter{
		server: engine,
		router: httpRouter.server.Group(apiServer.ApiPrefix),
	}
	return apiServer
}

func (apiServer *ApiServer) Register(apps ...App) {
	apiServer.locker.Lock()
	defer apiServer.locker.Unlock()
	for _, app := range apps {
		appName := app.Name()
		if _, ok := apiServer.apps[appName]; !ok {
			apiServer.apps[appName] = app
			rg := apiServer.rootRouter.Group(app.BaseRouter(), nil)
			app.Register(rg)
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
	err := apiServer.rootRouter.server.Run(apiServer.Addr)
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
