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
	"regexp"
	"strings"
	"sync"
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

	// Use define a API that for Controller Server Options for your Application.
	Use(option *ServerOption)
}

// IController represents a Controller(the MVC of Controller).
type IController interface {
	// Name define a API that returns Your controller name(such as resource name)
	// It's will be used as router prefix.
	Name() string
}

// ErrorHandler represents a http Error handler.
type ErrorHandler func(requestId string, httpRequest string, headers []string, stack string, err []*gin.Error)

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
	StartBeforeHandler     func(s *HostServer) error
	ShutDownBeforeHandler  func(s *HostServer) error
	BackendStoreHandler    func(cnf conf.Config) Store
	AppConfigHandler       func(cnf conf.BootStrapConfig) *conf.Config
	Crypto                 func(conf conf.Config) ICrypto
	SessionStateStore      func(conf conf.Config) ISessionStateManager
	AuthManager            IAuthManager
	PermissionManager      IPermissionManager
	StoreDbSetupHandler    StoreDbSetupHandler
	StoreCacheSetupHandler StoreCacheSetupHandler
	//ActionBeforeHandler    func(c *gin.Context)
	//ActionResponseHandler  func(c *gin.Context)
	cnf *conf.Config
	bcs *conf.BootStrapConfig
}

// HostServer represents a  Host Server.
type HostServer struct {
	name                string
	locker              sync.Mutex
	options             *ServerOption
	router              *Router
	apps                map[string]App
	conf                *conf.Config
	store               Store
	hash                ICryptoHash
	protect             ICryptoProtect
	authManager         IAuthManager
	sessionStateManager ISessionStateManager
	permissionManager   IPermissionManager
	validators          map[string]*regexp.Regexp
	httpErrHandlers     map[int][]ErrorHandler
	beforeHandler       []gin.HandlerFunc
	afterHandler        []gin.HandlerFunc
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
		AuthManager:            defaultAm,
		PermissionManager:      defaultPm,
		SessionStateStore: func(conf conf.Config) ISessionStateManager {
			return DefaultSessionStateManager(conf)
		},
		Crypto: func(conf conf.Config) ICrypto {
			c := conf.Service.Security.Crypto
			return DefaultCrypto(c.Protect.Secret, c.Hash.Salt)
		},
		bcs: bcs,
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
	server = &HostServer{
		options:         sopt,
		name:            sopt.Name,
		apps:            make(map[string]App),
		httpErrHandlers: make(map[int][]ErrorHandler),
		beforeHandler:   make([]gin.HandlerFunc, 0),
		afterHandler:    make([]gin.HandlerFunc, 0),
	}
	servers[sopt.Name] = server
	return server
}

// BeforeHooks register a global http handler API into the server, it's called AT handling http request before.
func (s *HostServer) BeforeHooks(handlers ...gin.HandlerFunc) {
	if len(handlers) == 0 {
		return
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	s.beforeHandler = append(s.beforeHandler, handlers...)
}

// BeforeHooks register a global http handler API into the server, it's called AT handled http request After.
func (s *HostServer) AfterHooks(handlers ...gin.HandlerFunc) {
	if len(handlers) == 0 {
		return
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	s.afterHandler = append(s.afterHandler, handlers...)
}

// HandleError register a global http error handler API into the server.
func (s *HostServer) HandleError(httpStatus int, handlers ...ErrorHandler) {
	if len(handlers) == 0 {
		return
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	h := s.httpErrHandlers[httpStatus]
	if len(h) == 0 {
		h = make([]ErrorHandler, len(handlers))
		copy(h, handlers)
	} else {
		h = append(h, handlers...)
	}
	s.httpErrHandlers[httpStatus] = h
}

// Register register a app instances into the server.
func (s *HostServer) Register(apps ...App) {
	for _, app := range apps {
		appName := app.Name()
		if _, ok := s.apps[appName]; !ok {
			s.apps[appName] = app
			// ServerOptions used AT first.
			app.Use(s.options)
		}
	}
}

// RegisterByPluginDir register app instances into the server with gw plugin mode.
func (s *HostServer) RegisterByPluginDir(dirs ...string) {
	for _, d := range dirs {
		rd, err := ioutil.ReadDir(d)
		if err != nil {
			logger.Error("read dir: %s", d)
			continue
		}
		for _, fi := range rd {
			pn := path.Join(d, fi.Name())
			if fi.IsDir() {
				s.RegisterByPluginDir(pn)
			} else {
				if !strings.HasSuffix(fi.Name(), s.options.PluginSymbolSuffix) {
					logger.Info("suffix not is %s, skipping file: %s", s.options.PluginSymbolSuffix, pn)
					continue
				}
				p, err := plugin.Open(pn)
				if err != nil {
					logger.Error("load plugin file: %s, err: %v", pn, err)
					continue
				}
				sym, err := p.Lookup(s.options.PluginSymbolName)
				if err != nil {
					logger.Error("file %s, err: %v", pn, err)
					continue
				}
				// TODO(Ocean): If symbol are pointer, how to conversions ?
				app, ok := sym.(App)
				if !ok {
					logger.Error("symbol %s in file %s did not is app.App interface.", s.options.PluginSymbolName, pn)
					continue
				}
				s.Register(app)
			}
		}
	}
}

func initial(s *HostServer) {
	//
	// Before Server start, Must initial all of Server Options
	// There are options may be changed by Custom app instance's .Use(...) APIs.
	cnf := s.options.AppConfigHandler(*s.options.bcs)
	crypto := s.options.Crypto(*cnf)

	s.options.cnf = cnf
	s.conf = cnf
	s.hash = crypto.Hash()
	s.protect = crypto.Protect()

	s.store = s.options.BackendStoreHandler(*cnf)
	s.sessionStateManager = s.options.SessionStateStore(*cnf)

	s.authManager = s.options.AuthManager
	s.permissionManager = s.options.PermissionManager

	registerValidators(s)

	// Gin Engine configure.
	gin.SetMode(s.options.Mode)
	g := gin.New()
	// g.Use(gin.Recovery())
	g.Use(globalState(s.options.Name))

	// Auth(login/logout) API routers.
	registerAuthRouter(cnf, s, g)

	var vars = make(map[string]string)
	vars["${PREFIX}"] = s.options.cnf.Service.Prefix
	g.Use(gwAuthChecker(vars, s.options.cnf.Service.Security.Auth.AllowUrls))

	if s.options.Mode == "debug" {
		g.Use(gin.Logger())
	}
	httpRouter := &Router{
		server: g,
	}

	// Must ensure store handler is not nil.
	if s.options.StoreDbSetupHandler == nil {
		s.options.StoreDbSetupHandler = appDefaultStoreDbSetupHandler
	}
	if s.options.StoreCacheSetupHandler == nil {
		s.options.StoreCacheSetupHandler = appDefaultStoreCacheSetupHandler
	}

	// initial routes.
	httpRouter.router = httpRouter.server.Group(s.options.Prefix)
	s.router = httpRouter
}

func registerValidators(s *HostServer) {

	p := s.conf.Service.Security.Auth
	s.validators = make(map[string]*regexp.Regexp)

	passportRegex, err := regexp.Compile(p.ParamPattern.Passport)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s for passport", p.ParamPattern.Passport))
	}
	secretRegex, err := regexp.Compile(p.ParamPattern.Secret)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s for secret", p.ParamPattern.Secret))
	}

	verifyCodeRegex, err := regexp.Compile(p.ParamPattern.VerifyCode)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s for verifyCode", p.ParamPattern.VerifyCode))
	}
	s.validators[p.ParamKey.Passport] = passportRegex
	s.validators[p.ParamKey.Secret] = secretRegex
	s.validators[p.ParamKey.VerifyCode] = verifyCodeRegex
}

func registerAuthRouter(cnf *conf.Config, server *HostServer, router *gin.Engine) {
	router.GET(cnf.Service.Security.Auth.LoginUrl, gwLogin)
	router.POST(cnf.Service.Security.Auth.LoginUrl, gwLogin)
	router.GET(cnf.Service.Security.Auth.LogoutUrl, gwLogout)
	router.POST(cnf.Service.Security.Auth.LogoutUrl, gwLogout)
}

func registerApps(server *HostServer) {
	for _, app := range server.apps {
		// app routers.
		logger.Info("register app: %s", app.Name())
		rg := server.router.Group(app.Router(), nil)
		app.Register(rg)

		// migrate
		logger.Info("migrate app: %s", app.Name())
		app.Migrate(server.store)
	}
}

// Serve represents start the Server.
func (s *HostServer) Serve() {
	//
	// All of app server initial AT here.
	initial(s)
	registerApps(s)

	// signal watch.
	sigs := make(chan os.Signal, 1)
	handler := s.options.StartBeforeHandler
	if handler != nil {
		err := s.options.StartBeforeHandler(s)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	logger.Info("Listening and serving HTTP on: %s", s.options.Addr)
	logger.ResetLogFormatter()
	var err error
	go func() {
		err = s.router.server.Run(s.options.Addr)
	}()
	// TODO(Ocean): has no better solution that can be waiting for gin.Serve() completed with non-block state?
	time.Sleep(time.Second * 1)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-sigs
	handler = s.options.ShutDownBeforeHandler
	if handler != nil {
		err := s.options.ShutDownBeforeHandler(s)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
	logger.Info("Shutdown server: %s, Addr: %s", s.options.Name, s.options.Addr)
}
