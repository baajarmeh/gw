package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
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

// ======================================== //
//                                          //
//    Define all of configuration Items     //
//                                          //
// ======================================== //

type Config struct {
	Service struct {
		Prefix   string `yaml:"prefix" toml:"prefix" json:"prefix"`
		Version  string `yaml:"version" toml:"version" json:"version"`
		Remarks  string `yaml:"remarks" toml:"remarks" json:"remarks"`
		Security struct {
			Auth struct {
				Disable     bool     `yaml:"disable" toml:"disable" json:"disable"`
				AllowedUrls []string `yaml:"allow-urls" toml:"allow-urls" json:"allow-urls"`
			}
		} `yaml:"security" toml:"security" json:"security"`
	} `yaml:"service" toml:"service" json:"service"`
	Common struct {
		Backend *Backend `yaml:"backend" toml:"backend" json:"backend"`
	} `yaml:"common" toml:"common" json:"common"`
}

type Backend struct {
	Db    *[]Db    `yaml:"db"`
	Cache *[]Cache `yaml:"cache"`
}

type Db struct {
	Driver   string            `yaml:"driver" toml:"driver" json:"driver"`
	Name     string            `yaml:"name" toml:"name" json:"name"`
	Addr     string            `yaml:"addr" toml:"addr" json:"addr"`
	Port     int               `yaml:"port" toml:"port" json:"port"`
	User     string            `yaml:"user" toml:"user" json:"user"`
	Password string            `yaml:"password" toml:"password" json:"password"`
	Database string            `yaml:"database" toml:"database" json:"database"`
	SSLMode  string            `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"`
	SSLCert  string            `yaml:"ssl_cert" toml:"ssl_cert" json:"ssl_cert"`
	Args     map[string]string `yaml:"args" toml:"args" json:"args"`
}

type Cache struct {
	Driver   string            `yaml:"driver" toml:"driver" json:"driver"`
	Name     string            `yaml:"name" toml:"name" json:"name"`
	Addr     string            `yaml:"addr" toml:"addr" json:"addr"`
	Port     int               `yaml:"port" toml:"port" json:"port"`
	User     string            `yaml:"user" toml:"user" json:"user"`
	Password string            `yaml:"password" toml:"password" json:"password"`
	Database string            `yaml:"database" toml:"database" json:"database"`
	SSLMode  string            `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"`
	SSLCert  string            `yaml:"ssl_cert" toml:"ssl_cert" json:"ssl_cert"`
	Args     map[string]string `yaml:"args" toml:"args" json:"args"`
}

// ============ End of configuration items ============= //

type ServerOption struct {
	Addr                  string
	Name                  string
	Mode                  string
	Restart               string
	ApiPrefix string
	StartBeforeHandler    func(server *ApiHostServer) error
	ShutDownBeforeHandler func(server *ApiHostServer) error
}

type ApiHostServer struct {
	locker                sync.Mutex
	router                *ApiRouter
	apps                  map[string]App
	option                *ServerOption
	conf                  *Config
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "oceanho.app"
	appDefaultRestart               = "always"
	appDefaultMode                  = "debug"
	appDefaultApiPrefix             = "api/v1"
	appDefaultApiVersion            = "Version 1.0"
	appDefaultStartBeforeHandler    = func(server *ApiHostServer) error { return nil }
	appDefaultShutDownBeforeHandler = func(server *ApiHostServer) error { return nil }
)

var (
	apiServers          map[string]*ApiHostServer
	apiServerSafeLocker sync.Mutex
)

func init() {
	apiServers = make(map[string]*ApiHostServer)
}

func NewServerOption() *ServerOption {
	conf := &ServerOption{
		Addr:                  appDefaultAddr,
		Name:                  appDefaultName,
		Restart:               appDefaultRestart,
		Mode:                  appDefaultMode,
		ApiPrefix: appDefaultApiPrefix,
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
	apiServerSafeLocker.Lock()
	defer apiServerSafeLocker.Unlock()
	apiServer, ok := apiServers[name]
	if ok {
		return apiServer
	}
	return nil
}

func New(conf *ServerOption) *ApiHostServer {
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
	apiServer = &ApiHostServer{
		apps:   make(map[string]App),
		option: conf,
	}
	httpRouter.router = httpRouter.server.Group(apiServer.option.ApiPrefix)
	apiServer.router = httpRouter
	apiServers[conf.Name] = apiServer
	return apiServer
}

func (apiServer *ApiHostServer) Register(apps ...App) {
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

func (apiServer *ApiHostServer) Serve() {
	sigs := make(chan os.Signal, 1)
	handler := apiServer.option.StartBeforeHandler
	if handler != nil {
		err := apiServer.option.StartBeforeHandler(apiServer)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	err := apiServer.router.server.Run(apiServer.option.Addr)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	handler = apiServer.option.ShutDownBeforeHandler
	if handler != nil {
		err := apiServer.option.ShutDownBeforeHandler(apiServer)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
}
