package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app/conf"
	_ "github.com/oceanho/gw/contrib/app/conf"
	"github.com/oceanho/gw/contrib/app/logger"
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
	StartBeforeHandler    func(server *ApiHostServer) error
	ShutDownBeforeHandler func(server *ApiHostServer) error
}

type ApiHostServer struct {
	locker sync.Mutex
	router *ApiRouter
	apps   map[string]App
	option *ServerOption
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
)

var (
	servers          map[string]*ApiHostServer
	serverSafeLocker sync.Mutex
)

func init() {
	servers = make(map[string]*ApiHostServer)
}

func NewServerOption() *ServerOption {
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
	}
	return conf
}

func Default() *ApiHostServer {
	return New(NewServerOption())
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

func New(conf *ServerOption) *ApiHostServer {
	serverSafeLocker.Lock()
	defer serverSafeLocker.Unlock()
	server, ok := servers[conf.Name]
	if ok {
		return server
	}
	gin.SetMode(conf.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	if conf.Mode == "debug" {
		engine.Use(gin.Logger())
	}
	httpRouter := &ApiRouter{
		server: engine,
	}
	server = &ApiHostServer{
		apps:   make(map[string]App),
		option: conf,
	}
	httpRouter.router = httpRouter.server.Group(server.option.ApiPrefix)
	server.router = httpRouter
	servers[conf.Name] = server
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
				if !strings.HasSuffix(fi.Name(), server.option.PluginSymbolSuffix) {
					logger.Info("suffix not is %s, skipping file: %s, err: %v", server.option.PluginSymbolSuffix, pn)
					continue
				}
				p, err := plugin.Open(pn)
				if err != nil {
					logger.Error("load plugin file: %s, err: %v", pn, err)
					continue
				}
				sym, err := p.Lookup(server.option.PluginSymbolName)
				if err != nil {
					logger.Error("file %s, err: %v", pn, err)
					continue
				}
				app, ok := sym.(App)
				if !ok {
					logger.Error("symbol %s in file %s did not is app.App interface.", server.option.PluginSymbolName, pn)
					continue
				}
				server.Register(app)
			}
		}
	}
}

func (server *ApiHostServer) Serve() {
	sigs := make(chan os.Signal, 1)
	handler := server.option.StartBeforeHandler
	if handler != nil {
		err := server.option.StartBeforeHandler(server)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	err := server.router.server.Run(server.option.Addr)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	handler = server.option.ShutDownBeforeHandler
	if handler != nil {
		err := server.option.ShutDownBeforeHandler(server)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
}
