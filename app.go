package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"plugin"
	"strings"
	"syscall"
	"time"
)

// App represents a application.
type App interface {

	// Name define a API that return as your app name.
	Name() string

	// Router define a API that should be return your app base route path.
	// It's will be used to create a new *RouteGroup object and that used by Register(...) .
	Router() string

	// Register define a API that for register your app router inside.
	Register(router *RouterGroup)

	// Migrate define a API that for create your app database migrations inside.
	Migrate(store Store)
}

// IController represents a Controller(the MVC of Controller).
type IController interface {
	// Name define a API that returns Your controller name(such as resource name)
	// It's will be used as router prefix.
	Name() string

	// Create define a API
	Create(ctx *Context) IController

	// OnDestroy define a API
	OnDestroy(ctx *Context) error
}

// ServerOption represents a Server Options.
type ServerOption struct {
	Addr                   string
	Name                   string
	Mode                   string
	Restart                string
	Prefix                 string
	PluginDir              string
	PluginSymbolName       string
	PluginSymbolSuffix     string
	StartBeforeHandler     func(server *HostServer) error
	ShutDownBeforeHandler  func(server *HostServer) error
	BackendStoreHandler    func(cnf conf.Config) Store
	AppConfigHandler       func(cnf conf.BootStrapConfig) *conf.Config
	StoreDbSetupHandler    StoreDbSetupHandler
	StoreCacheSetupHandler StoreCacheSetupHandler
	cnf                    *conf.Config
	bcs                    *conf.BootStrapConfig
}

// HostServer represents a  Host Server.
type HostServer struct {
	options *ServerOption
	router  *Router
	apps    map[string]App
	conf    *conf.Config
	store   Store
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "gw.app"
	appDefaultRestart               = "always"
	appDefaultMode                  = "debug"
	appDefaultPrefix                = "/api/v1"
	appDefaultPluginSymbolName      = "AppPlugin"
	appDefaultPluginSymbolSuffix    = ".so"
	appDefaultVersion               = "Version 1.0"
	appDefaultStartBeforeHandler    = func(server *HostServer) error { return nil }
	appDefaultShutDownBeforeHandler = func(server *HostServer) error { return nil }
	appDefaultBackendHandler        = func(cnf conf.Config) Store {
		return DefaultBackend(cnf)
	}
	appDefaultAppConfigHandler = func(cnf conf.BootStrapConfig) *conf.Config {
		return conf.NewConfigByBootStrapConfig(&cnf)
	}
	appDefaultStoreDbSetupHandler = func(c gin.Context, db *gorm.DB, user User) *gorm.DB {
		return db
	}
	appDefaultStoreCacheSetupHandler = func(c gin.Context, client *redis.Client, user User) *redis.Client {
		return client
	}
	internLogFormatter = "[$prefix-$level] - $msg\n"
)

var (
	servers map[string]*HostServer
)

func init() {
	servers = make(map[string]*HostServer)
	logger.SetLogFormatter(internLogFormatter)
}

// NewServerOption returns a *ServerOption with bcs.
func NewServerOption(bcs *conf.BootStrapConfig) *ServerOption {
	mode := os.Getenv(gin.EnvGinMode)
	if mode == "" {
		mode = appDefaultMode
	}
	conf := &ServerOption{
		Mode:                   mode,
		Addr:                   appDefaultAddr,
		Name:                   appDefaultName,
		Restart:                appDefaultRestart,
		Prefix:                 appDefaultPrefix,
		AppConfigHandler:       appDefaultAppConfigHandler,
		PluginSymbolName:       appDefaultPluginSymbolName,
		PluginSymbolSuffix:     appDefaultPluginSymbolSuffix,
		StartBeforeHandler:     appDefaultStartBeforeHandler,
		ShutDownBeforeHandler:  appDefaultShutDownBeforeHandler,
		BackendStoreHandler:    appDefaultBackendHandler,
		StoreDbSetupHandler:    appDefaultStoreDbSetupHandler,
		StoreCacheSetupHandler: appDefaultStoreCacheSetupHandler,
		bcs:                    bcs,
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
	server, ok := servers[name]
	if ok {
		return server
	}
	return nil
}

// New return a  Server with ServerOptions.
func New(sopt *ServerOption) *HostServer {
	server, ok := servers[sopt.Name]
	if ok {
		logger.Warn("duplicated server, name: %s", sopt.Name)
		return server
	}
	gin.SetMode(sopt.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gwState(sopt.Name))
	if sopt.Mode == "debug" {
		engine.Use(gin.Logger())
	}
	httpRouter := &Router{
		server: engine,
	}
	cnf := sopt.AppConfigHandler(*sopt.bcs)
	sopt.cnf = cnf
	server = &HostServer{
		options: sopt,
		apps:    make(map[string]App),
		conf:    cnf,
		store:   sopt.BackendStoreHandler(*cnf),
	}
	// Must ensure store handler is not nil.
	if server.options.StoreDbSetupHandler == nil {
		server.options.StoreDbSetupHandler = appDefaultStoreDbSetupHandler
	}
	if server.options.StoreCacheSetupHandler == nil {
		server.options.StoreCacheSetupHandler = appDefaultStoreCacheSetupHandler
	}

	// initial routes.
	httpRouter.router = httpRouter.server.Group(server.options.Prefix)
	server.router = httpRouter
	servers[sopt.Name] = server
	return server
}

// Register register app instances into the server.
func (server *HostServer) Register(apps ...App) {
	for _, app := range apps {
		appName := app.Name()
		if _, ok := server.apps[appName]; !ok {
			server.apps[appName] = app
			logger.Info("register app: %s", appName)
			rg := server.router.Group(app.Router(), nil)
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
				if !strings.HasSuffix(fi.Name(), server.options.PluginSymbolSuffix) {
					logger.Info("suffix not is %s, skipping file: %s", server.options.PluginSymbolSuffix, pn)
					continue
				}
				p, err := plugin.Open(pn)
				if err != nil {
					logger.Error("load plugin file: %s, err: %v", pn, err)
					continue
				}
				sym, err := p.Lookup(server.options.PluginSymbolName)
				if err != nil {
					logger.Error("file %s, err: %v", pn, err)
					continue
				}
				// TODO(Ocean): If symbol are pointer, how to conversions ?
				app, ok := sym.(App)
				if !ok {
					logger.Error("symbol %s in file %s did not is app.App interface.", server.options.PluginSymbolName, pn)
					continue
				}
				server.Register(app)
			}
		}
	}
}

// Serve represents start the Server.
func (server *HostServer) Serve() {
	// before server starting. Try migrates for all registered Apps.
	for _, p := range server.apps {
		logger.Info("Migrate app: %s", p.Name())
		p.Migrate(server.store)
	}

	// signal watch.
	sigs := make(chan os.Signal, 1)
	handler := server.options.StartBeforeHandler
	if handler != nil {
		err := server.options.StartBeforeHandler(server)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	logger.Info("Listening and serving HTTP on: %s", server.options.Addr)
	logger.ResetLogFormatter()
	var err error
	go func() {
		err = server.router.server.Run(server.options.Addr)
	}()
	// TODO(Ocean): has no better solution that can be waiting for gin.Serve() completed with non-block state?
	time.Sleep(time.Second * 1)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-sigs
	handler = server.options.ShutDownBeforeHandler
	if handler != nil {
		err := server.options.ShutDownBeforeHandler(server)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
	logger.Info("Shutdown server: %s, Addr: %s", server.options.Name, server.options.Addr)
}
