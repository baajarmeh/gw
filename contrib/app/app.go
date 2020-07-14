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

// ======================================== //
//                                          //
//    Define all of configuration Items     //
//                                          //
// ======================================== //

type Config struct {
	Api    Api `yaml:"api"`
	Common struct {
		Backend *Backend `yaml:"backend"`
	} `yaml:"common"`
}

type Api struct {
	Prefix   string `yaml:"prefix"`
	Version  string `yaml:"version"`
	Remarks  string `yaml:"version"`
	Security struct {
		Auth struct {
			Disable     bool     `yaml:"disable"`
			AllowedUrls []string `yaml:"allow-urls"`
		}
	} `yaml:"security"`
}

type Backend struct {
	Db    *[]Db    `yaml:"db"`
	Cache *[]Cache `yaml:"cache"`
}

type Db struct {
	Driver   string            `yaml:"driver"`
	Name     string            `yaml:"name"`
	Addr     string            `yaml:"addr"`
	Port     int               `yaml:"port"`
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	Database string            `yaml:"database"`
	SSLMode  string            `yaml:"ssl_mode"`
	SSLCert  string            `yaml:"ssl_cert"`
	Args     map[string]string `yaml:"args"`
}

type Cache struct {
	Name     string            `yaml:"name"`
	Addr     string            `yaml:"addr"`
	Port     int               `yaml:"port"`
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	Database int               `yaml:"database"`
	SSLMode  string            `yaml:"ssl_mode"`
	SSLCert  string            `yaml:"ssl_cert"`
	Args     map[string]string `yaml:"args"`
}

// ============ End of configuration items ============= //

type Option struct {
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
	apps                  map[string]App
	option                *Option
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

var (
	apiServers          map[string]*ApiServer
	apiServerSafeLocker sync.Mutex
)

func init() {
	apiServers = make(map[string]*ApiServer)
}

func NewOption() *Option {
	conf := &Option{
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
	return New(NewOption())
}

func GetDefaultApiServer() *ApiServer {
	return GetApiServer(appDefaultName)
}

func GetApiServer(name string) *ApiServer {
	apiServerSafeLocker.Lock()
	defer apiServerSafeLocker.Unlock()
	apiServer, ok := apiServers[name]
	if ok {
		return apiServer
	}
	return nil
}

func New(conf *Option) *ApiServer {
	apiServerSafeLocker.Lock()
	defer apiServerSafeLocker.Unlock()
	apiServer, ok := apiServers[conf.Name]
	if ok {
		return apiServer
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
	apiServer = &ApiServer{
		Addr:                  conf.Addr,
		Name:                  conf.Name,
		Restart:               conf.Restart,
		Mode:                  conf.Mode,
		StartBeforeHandler:    conf.StartBeforeHandler,
		ShutDownBeforeHandler: conf.ShutDownBeforeHandler,
		ApiPrefix:             conf.ApiPrefix,
		ApiVersion:            conf.ApiVersion,
		apps:                  make(map[string]App),
		option:                conf,
	}
	httpRouter.router = httpRouter.server.Group(apiServer.ApiPrefix)
	apiServer.router = httpRouter
	apiServers[conf.Name] = apiServer
	return apiServer
}

func (apiServer *ApiServer) Register(apps ...App) {
	apiServer.locker.Lock()
	defer apiServer.locker.Unlock()
	for _, app := range apps {
		appName := app.Name()
		if _, ok := apiServer.apps[appName]; !ok {
			apiServer.apps[appName] = app
			fmt.Printf("\n[Ocean GW]  register app: %s\n", appName)
			rg := apiServer.router.Group(app.BaseRouter(), nil)
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
