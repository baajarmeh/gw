package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app/conf"
	_ "github.com/oceanho/gw/contrib/app/conf"
	"github.com/oceanho/gw/contrib/app/logger"
	"github.com/oceanho/gw/contrib/app/store"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"plugin"
	"strings"
	"sync"
	"syscall"
)

// App represents a application, Should be implement this.
type App interface {

	// Name define a API that return as your app name.
	Name() string

	// BaseRouter define a API that should be return your app base route path.
	// It's will be to create a new *ApiRouteGroup object and that used by Register(...) API.
	BaseRouter() string

	// Register define a API that for register your app router inside.
	Register(router *ApiRouteGroup)
}

type ServerOption struct {
	Addr                  string
	Name                  string
	Mode                  string
	Restart               string
	ApiPrefix             string
	PluginDir             string
	PluginSymbolName      string
	PluginSymbolSuffix    string
	bcs                   *conf.BootStrapConfig
	StartBeforeHandler    func(server *ApiHostServer) error
	ShutDownBeforeHandler func(server *ApiHostServer) error
	BackendStoreHandler   func(cnf conf.Config) store.Backend
	AppConfigHandler      func(cnf conf.BootStrapConfig) *conf.Config
}

type ApiHostServer struct {
	Options *ServerOption
	// private properties.
	locker sync.Mutex
	router *ApiRouter
	apps   map[string]App
	conf   *conf.Config
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "oceanho.app"
	appDefaultRestart               = "always"
	appDefaultMode                  = "debug"
	appDefaultApiPrefix             = "api/v1"
	appDefaultPluginSymbolName      = "AppPlugin"
	appDefaultPluginSymbolSuffix    = ".so"
	appDefaultApiVersion            = "Version 1.0"
	appDefaultStartBeforeHandler    = func(server *ApiHostServer) error { return nil }
	appDefaultShutDownBeforeHandler = func(server *ApiHostServer) error { return nil }
	appDefaultBackendHandler        = func(cnf conf.Config) store.Backend {
		return store.Default(cnf)
	}
	appAppConfigHandler = func(cnf conf.BootStrapConfig) *conf.Config {
		return conf.NewConfigByBootStrapConfig(&cnf)
	}
)

var (
	servers          map[string]*ApiHostServer
	serverSafeLocker sync.Mutex
)

func init() {
	servers = make(map[string]*ApiHostServer)
}

func NewServerOption(bcs *conf.BootStrapConfig) *ServerOption {
	conf := &ServerOption{
		Addr:                  appDefaultAddr,
		Name:                  appDefaultName,
		Restart:               appDefaultRestart,
		Mode:                  appDefaultMode,
		ApiPrefix:             appDefaultApiPrefix,
		PluginSymbolName:      appDefaultPluginSymbolName,
		PluginSymbolSuffix:    appDefaultPluginSymbolSuffix,
		StartBeforeHandler:    appDefaultStartBeforeHandler,
		ShutDownBeforeHandler: appDefaultShutDownBeforeHandler,
		BackendStoreHandler:   appDefaultBackendHandler,
		AppConfigHandler:      appAppConfigHandler,
		bcs:                   bcs,
	}
	return conf
}

func Default() *ApiHostServer {
	bcs := conf.DefaultBootStrapConfig()
	return New(NewServerOption(bcs))
}

func GetDefaultApiServer() *ApiHostServer {
	return GetApiServer(appDefaultName)
}

func GetApiServer(name string) *ApiHostServer {
	serverSafeLocker.Lock()
	defer serverSafeLocker.Unlock()
	server, ok := servers[name]
	if ok {
		return server
	}
	return nil
}

func New(sopt *ServerOption) *ApiHostServer {
	serverSafeLocker.Lock()
	defer serverSafeLocker.Unlock()
	server, ok := servers[sopt.Name]
	if ok {
		logger.Warn("duplicated server, name: %s", sopt.Name)
		return server
	}
	gin.SetMode(sopt.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	if sopt.Mode == "debug" {
		engine.Use(gin.Logger())
	}
	httpRouter := &ApiRouter{
		server: engine,
	}
	server = &ApiHostServer{
		apps:    make(map[string]App),
		Options: sopt,
		conf:    appAppConfigHandler(*sopt.bcs),
	}
	//
	// Initial all of internal components AT here.
	appConf := *server.conf
	// backend initialization
	store.Initial(appConf, server.Options.BackendStoreHandler)

	// initial routes.
	httpRouter.router = httpRouter.server.Group(server.Options.ApiPrefix)
	server.router = httpRouter
	servers[sopt.Name] = server
	return server
}

func (server *ApiHostServer) Register(apps ...App) {
	server.locker.Lock()
	defer server.locker.Unlock()
	for _, app := range apps {
		appName := app.Name()
		if _, ok := server.apps[appName]; !ok {
			server.apps[appName] = app
			logger.Info("register app: %s", appName)
			rg := server.router.Group(app.BaseRouter(), nil)
			app.Register(rg)
		}
	}
}

func (server *ApiHostServer) RegisterByPluginDir(dirs ...string) {
	for _, d := range dirs {
		rd, err := ioutil.ReadDir(d)
		if err != nil {
			logger.Error("read dir: %s", d)
			continue
		}
		for _, fi := range rd {
			pn := path.Join(d, fi.Name())
			if fi.IsDir() {
				server.RegisterByPluginDir(pn)
			} else {
				if !strings.HasSuffix(fi.Name(), server.Options.PluginSymbolSuffix) {
					logger.Info("suffix not is %s, skipping file: %s", server.Options.PluginSymbolSuffix, pn)
					continue
				}
				p, err := plugin.Open(pn)
				if err != nil {
					logger.Error("load plugin file: %s, err: %v", pn, err)
					continue
				}
				sym, err := p.Lookup(server.Options.PluginSymbolName)
				if err != nil {
					logger.Error("file %s, err: %v", pn, err)
					continue
				}
				// TODO(Ocean): If symbol are pointer, how to conversions ?
				app, ok := sym.(App)
				if !ok {
					logger.Error("symbol %s in file %s did not is app.App interface.", server.Options.PluginSymbolName, pn)
					continue
				}
				server.Register(app)
			}
		}
	}
}

func (server *ApiHostServer) Serve() {
	sigs := make(chan os.Signal, 1)
	handler := server.Options.StartBeforeHandler
	if handler != nil {
		err := server.Options.StartBeforeHandler(server)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	err := server.router.server.Run(server.Options.Addr)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	handler = server.Options.ShutDownBeforeHandler
	if handler != nil {
		err := server.Options.ShutDownBeforeHandler(server)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
}
