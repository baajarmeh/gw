package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
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

	// Name define a  that return as your app name.
	Name() string

	// BaseRouter define a  that should be return your app base route path.
	// It's will be to create a new *RouteGroup object and that used by Register(...) .
	BaseRouter() string

	// Register define a  that for register your app router inside.
	Register(router *RouteGroup)
}

// ServerOption represents a Server Options.
type ServerOption struct {
	Addr                  string
	Name                  string
	Mode                  string
	Restart               string
	Prefix                string
	PluginDir             string
	PluginSymbolName      string
	PluginSymbolSuffix    string
	bcs                   *conf.BootStrapConfig
	StartBeforeHandler    func(server *HostServer) error
	ShutDownBeforeHandler func(server *HostServer) error
	BackendStoreHandler   func(cnf conf.Config) Store
	AppConfigHandler      func(cnf conf.BootStrapConfig) *conf.Config
}

// HostServer represents a  Host Server.
type HostServer struct {
	Options *ServerOption
	// private properties.
	locker sync.Mutex
	router *Router
	apps   map[string]App
	conf   *conf.Config
	store  Store
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "oceanho.app"
	appDefaultRestart               = "always"
	appDefaultMode                  = "debug"
	appDefaultPrefix                = "/v1"
	appDefaultPluginSymbolName      = "AppPlugin"
	appDefaultPluginSymbolSuffix    = ".so"
	appDefaultVersion               = "Version 1.0"
	appDefaultStartBeforeHandler    = func(server *HostServer) error { return nil }
	appDefaultShutDownBeforeHandler = func(server *HostServer) error { return nil }
	appDefaultBackendHandler        = func(cnf conf.Config) Store {
		return DefaultBackend(cnf)
	}
	appAppConfigHandler = func(cnf conf.BootStrapConfig) *conf.Config {
		return conf.NewConfigByBootStrapConfig(&cnf)
	}
)

var (
	servers          map[string]*HostServer
	serverSafeLocker sync.Mutex
)

func init() {
	servers = make(map[string]*HostServer)
}

// NewServerOption returns a *ServerOption with bcs.
func NewServerOption(bcs *conf.BootStrapConfig) *ServerOption {
	conf := &ServerOption{
		Addr:                  appDefaultAddr,
		Name:                  appDefaultName,
		Restart:               appDefaultRestart,
		Mode:                  appDefaultMode,
		Prefix:                appDefaultPrefix,
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

// Default returns a default HostServer(the server instance's bcs,svr are default config items.)
func Default() *HostServer {
	bcs := conf.DefaultBootStrapConfig()
	return New(NewServerOption(bcs))
}

// GetDefaultHostServer return a default  Server of has registered.
func GetDefaultHostServer() *HostServer {
	return GetGwServer(appDefaultName)
}

// GetGwServer return a default  Server by name of has registered.
func GetGwServer(name string) *HostServer {
	serverSafeLocker.Lock()
	defer serverSafeLocker.Unlock()
	server, ok := servers[name]
	if ok {
		return server
	}
	return nil
}

// New return a  Server with ServerOptions.
func New(sopt *ServerOption) *HostServer {
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
	httpRouter := &Router{
		server: engine,
	}
	server = &HostServer{
		apps:    make(map[string]App),
		Options: sopt,
		conf:    appAppConfigHandler(*sopt.bcs),
	}
	//
	// Initial all of internal components AT here.
	appConf := *server.conf
	// backend initialization
	Initial(appConf, server.Options.BackendStoreHandler)

	// initial routes.
	httpRouter.router = httpRouter.server.Group(server.Options.Prefix)
	server.router = httpRouter
	servers[sopt.Name] = server
	return server
}

// Register register app instances into the server.
func (server *HostServer) Register(apps ...App) {
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

// RegisterByPluginDir register app instances into the server with gw plugin mode.
func (server *HostServer) RegisterByPluginDir(dirs ...string) {
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

// Serve represents start the Server.
func (server *HostServer) Serve() {
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
