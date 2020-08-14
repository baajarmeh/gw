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

	// Migrate define a API that for create your app's database migrations and permission initialization inside.
	Migrate(ctx MigrationContext)

	// Use define a API that for RestAPI Server Options for your Application.
	Use(option *ServerOption)
}

// MigrationContext represents a Migration Context Object.
type MigrationContext struct {
	Store             Store
	PermissionManager IPermissionManager
}

// IRestAPI represents a Rest Style API instance.
type IRestAPI interface {
	// Name define a API that returns Your RestAPI name(such as resource name)
	// It's will be used as router prefix.
	Name() string
}

// ErrorHandler represents a http Error handler.
type ErrorHandler func(requestId string, httpRequest string, headers []string, stack string, err []*gin.Error)

// ServerOption represents a Server Options.
type ServerOption struct {
	Addr                     string
	Name                     string
	Restart                  string
	Prefix                   string
	PluginDir                string
	PluginSymbolName         string
	PluginSymbolSuffix       string
	StartBeforeHandler       func(s *HostServer) error
	ShutDownBeforeHandler    func(s *HostServer) error
	BackendStoreHandler      func(cnf conf.ApplicationConfig) Store
	AppConfigHandler         func(cnf conf.BootConfig) *conf.ApplicationConfig
	Crypto                   func(conf conf.ApplicationConfig) ICrypto
	AuthManager              IAuthManager
	PermissionManagerHandler PermissionManagerHandler
	StoreDbSetupHandler      StoreDbSetupHandler
	SessionStateManager      SessionStateHandler
	StoreCacheSetupHandler   StoreCacheSetupHandler
	RespBodyBuildFunc        RespBodyCreationHandler
	cnf                      *conf.ApplicationConfig
	bcs                      *conf.BootConfig
}

// HostServer represents a  Host Server.
type HostServer struct {
	Store               Store
	Hash                ICryptoHash
	Protect             ICryptoProtect
	AuthManager         IAuthManager
	SessionStateManager ISessionStateManager
	PermissionManager   IPermissionManager
	RespBodyBuildFunc   RespBodyCreationHandler
	// private
	state                  int
	locker                 sync.Mutex
	options                *ServerOption
	router                 *Router
	apps                   map[string]App
	conf                   *conf.ApplicationConfig
	httpErrHandlers        map[int][]ErrorHandler
	hooks                  []*Hook
	beforeHooks            []*Hook
	afterHooks             []*Hook
	afterHookMaxIdx        int
	authParamValidators    map[string]*regexp.Regexp
	storeDbSetupHandler    StoreDbSetupHandler
	storeCacheSetupHandler StoreCacheSetupHandler
	migrateContext         MigrationContext
	serverExitSignal       chan struct{}
	serverStartDone        chan struct{}
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "gw.app"
	appDefaultRestart               = "always"
	appDefaultPrefix                = "/api/v1"
	appDefaultPluginSymbolName      = "AppPlugin"
	appDefaultPluginSymbolSuffix    = ".so"
	appDefaultStartBeforeHandler    = func(server *HostServer) error { return nil }
	appDefaultShutdownBeforeHandler = func(server *HostServer) error { return nil }
	appDefaultBackendHandler        = func(cnf conf.ApplicationConfig) Store {
		return DefaultBackend(cnf)
	}
	appDefaultAppConfigHandler = func(cnf conf.BootConfig) *conf.ApplicationConfig {
		return conf.NewConfigWithBootConfig(&cnf)
	}
	appDefaultStoreDbSetupHandler = func(c Context, db *gorm.DB) *gorm.DB {
		return db
	}
	appDefaultStoreCacheSetupHandler = func(c Context, client *redis.Client, user User) *redis.Client {
		return client
	}
	internLogFormatter = "[$prefix-$level] $msg\n"
)

var (
	servers map[string]*HostServer
)

func init() {
	servers = make(map[string]*HostServer)
	logger.SetLogFormatter(internLogFormatter)
}

// NewServerOption returns a *ServerOption with bcs.
func NewServerOption(bcs *conf.BootConfig) *ServerOption {
	conf := &ServerOption{
		Addr:                   appDefaultAddr,
		Name:                   appDefaultName,
		Restart:                appDefaultRestart,
		Prefix:                 appDefaultPrefix,
		AppConfigHandler:       appDefaultAppConfigHandler,
		PluginSymbolName:       appDefaultPluginSymbolName,
		PluginSymbolSuffix:     appDefaultPluginSymbolSuffix,
		StartBeforeHandler:     appDefaultStartBeforeHandler,
		ShutDownBeforeHandler:  appDefaultShutdownBeforeHandler,
		BackendStoreHandler:    appDefaultBackendHandler,
		StoreDbSetupHandler:    appDefaultStoreDbSetupHandler,
		StoreCacheSetupHandler: appDefaultStoreCacheSetupHandler,
		AuthManager:            defaultAm,
		PermissionManagerHandler: func(conf conf.ApplicationConfig, store Store) IPermissionManager {
			return DefaultPermissionManager(conf, store)
		},
		SessionStateManager: func(conf conf.ApplicationConfig) ISessionStateManager {
			return DefaultSessionStateManager(conf)
		},
		Crypto: func(conf conf.ApplicationConfig) ICrypto {
			c := conf.Security.Crypto
			return DefaultCrypto(c.Protect.Secret, c.Hash.Salt)
		},
		RespBodyBuildFunc: defaultRespBodyBuildFunc,
		bcs:               bcs,
	}
	return conf
}

// DefaultServer returns a default HostServer(the server instance's bcs,svr are default config items.)
func DefaultServer() *HostServer {
	bcs := conf.DefaultBootConfig()
	return NewServerWithOption(NewServerOption(bcs))
}

// NewServerWithName returns a specifies name, other use default HostServer(the server instance's bcs,svr are default config items.)
func NewServerWithName(name string) *HostServer {
	bcs := conf.DefaultBootConfig()
	opts := NewServerOption(bcs)
	opts.Name = name
	return NewServerWithOption(opts)
}

// NewServerWithNameAddr returns a specifies name,addr, other use default HostServer(the server instance's bcs,svr are default config items.)
func NewServerWithNameAddr(name, addr string) *HostServer {
	bcs := conf.DefaultBootConfig()
	opts := NewServerOption(bcs)
	opts.Name = name
	opts.Addr = addr
	return NewServerWithOption(opts)
}

// NewServerWithName returns a specifies addr, other use default HostServer(the server instance's bcs,svr are default config items.)
func NewServerWithAddr(addr string) *HostServer {
	bcs := conf.DefaultBootConfig()
	opts := NewServerOption(bcs)
	opts.Addr = addr
	return NewServerWithOption(opts)
}

// New return a  Server with ServerOptions.
func NewServerWithOption(sopt *ServerOption) *HostServer {
	server, ok := servers[sopt.Name]
	if ok {
		logger.Warn("duplicated server, name: %s", sopt.Name)
		return server
	}
	server = &HostServer{
		options:             sopt,
		hooks:               make([]*Hook, 0),
		beforeHooks:         make([]*Hook, 0),
		afterHooks:          make([]*Hook, 0),
		apps:                make(map[string]App),
		httpErrHandlers:     make(map[int][]ErrorHandler),
		authParamValidators: make(map[string]*regexp.Regexp),
		serverExitSignal:    make(chan struct{}, 1),
		serverStartDone:     make(chan struct{}, 1),
	}
	servers[sopt.Name] = server
	return server
}

// AddHook register a global http handler API into the server.
func (s *HostServer) AddHook(handlers ...*Hook) {
	if len(handlers) == 0 {
		return
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	for j := 0; j < len(handlers); j++ {
		for i := 0; i < len(s.hooks); i++ {
			if s.hooks[i].Name == handlers[j].Name {
				s.hooks[i] = handlers[j]
				continue
			}
			s.hooks = append(s.hooks, handlers[j])
		}
	}
}

// DeleteHook delete a global hook from has registered before hooks.
func (s *HostServer) DeleteHook(name string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for i := 0; i < len(s.hooks); i++ {
		if s.hooks[i].Name == name {
			left := s.hooks[0:i]
			right := s.hooks[i+1:]
			s.hooks = left
			s.hooks = append(s.hooks, right...)
			break
		}
	}
}

// ReplaceHook replace a global hook from has registered before hooks.
func (s *HostServer) ReplaceHook(name string, hookHandler *Hook) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for i := 0; i < len(s.hooks); i++ {
		if s.hooks[i].Name == name {
			s.hooks[i] = hookHandler
			break
		}
	}
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
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, app := range apps {
		appName := app.Name()
		if _, ok := s.apps[appName]; !ok {
			s.apps[appName] = app
			// ServerOptions used AT first.
			app.Use(s.options)
		}
	}
}

// Patch patch up for HostServer with apps.
// It's Usage scenarios are need change server Options with other ge App but not need that's router features.
//
// Example:
// Your have two gw App
// 1. UAP(your app one), It's implementation User/Permission features
// 2. DevOps(your app two), It's implementation DevOps features
//    It's dependencies on UAP, but the UAP's router not necessary on DevOps service.
//
// Now, We can separate deployment of UAP and DevOps
// And we should be patch up UAP into DevOps app.
func (s *HostServer) Patch(apps ...App) {
	if len(apps) == 0 {
		return
	}
	ctx := MigrationContext{
		Store:             s.Store,
		PermissionManager: s.PermissionManager,
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, app := range apps {
		appName := app.Name()
		if _, ok := s.apps[appName]; !ok {
			app.Use(s.options)
			app.Migrate(ctx)
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
	s.Hash = crypto.Hash()
	s.Protect = crypto.Protect()

	s.Store = s.options.BackendStoreHandler(*cnf)
	s.SessionStateManager = s.options.SessionStateManager(*cnf)

	s.AuthManager = s.options.AuthManager
	s.PermissionManager = s.options.PermissionManagerHandler(*cnf, s.Store)

	registerAuthParamValidators(s)

	// gin engine.
	g := gin.New()
	// g.Use(gin.Recovery())
	g.Use(gwState(s.options.Name))

	// Auth(login/logout) API routers.
	registerAuthRouter(cnf, g)

	// global Auth middleware.
	g.Use(gwAuthChecker(s.options.cnf.Security.Auth.AllowUrls))

	if gin.IsDebugging() {
		g.Use(gin.Logger())
	}

	httpRouter := &Router{
		server:      g,
		prefix:      s.options.Prefix,
		routerInfos: make([]RouterInfo, 0),
	}

	// Must ensure store handler is not nil.
	if s.options.StoreDbSetupHandler == nil {
		s.options.StoreDbSetupHandler = appDefaultStoreDbSetupHandler
	}
	if s.options.StoreCacheSetupHandler == nil {
		s.options.StoreCacheSetupHandler = appDefaultStoreCacheSetupHandler
	}
	if s.options.RespBodyBuildFunc == nil {
		s.options.RespBodyBuildFunc = s.RespBodyBuildFunc
	}
	if s.storeDbSetupHandler == nil {
		s.storeDbSetupHandler = s.options.StoreDbSetupHandler
	}
	if s.storeCacheSetupHandler == nil {
		s.storeCacheSetupHandler = s.options.StoreCacheSetupHandler
	}
	if s.RespBodyBuildFunc == nil {
		s.RespBodyBuildFunc = s.options.RespBodyBuildFunc
	}
	// initial routes.
	httpRouter.router = httpRouter.server.Group(s.options.Prefix)
	s.router = httpRouter

	if s.conf.Server.ListenAddr != "" && s.options.Addr == appDefaultAddr {
		s.options.Addr = s.conf.Server.ListenAddr
	}
	// permission manager initial
	s.PermissionManager.Initial()
}

func registerAuthParamValidators(s *HostServer) {
	p := s.conf.Security.Auth
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
	s.authParamValidators[p.ParamKey.Passport] = passportRegex
	s.authParamValidators[p.ParamKey.Secret] = secretRegex
	s.authParamValidators[p.ParamKey.VerifyCode] = verifyCodeRegex
}

func registerAuthRouter(cnf *conf.ApplicationConfig, router *gin.Engine) {
	authServer := cnf.Security.AuthServer
	if authServer.EnableAuthServe {
		for _, m := range authServer.LogIn.Methods {
			router.Handle(strings.ToUpper(m), authServer.LogIn.Url, gwLogin)
		}
		for _, m := range authServer.LogOut.Methods {
			router.Handle(strings.ToUpper(m), authServer.LogOut.Url, gwLogout)
		}
	}
}

func registerApps(s *HostServer) {
	ctx := MigrationContext{
		Store:             s.Store,
		PermissionManager: s.PermissionManager,
	}
	for _, app := range s.apps {
		//
		// app routers.
		logger.Info("register app: %s", app.Name())
		rg := s.router.Group(app.Router(), nil)
		app.Register(rg)

		// migrate
		logger.Info("migrate app: %s", app.Name())
		app.Migrate(ctx)
	}
}

func prepareHooks(s *HostServer) {
	if len(s.hooks) < 1 {
		return
	}
	for i := 0; i < len(s.hooks); i++ {
		if s.hooks[i].OnAfter != nil {
			s.afterHooks = append(s.afterHooks, s.hooks[i])
		}
		if s.hooks[i].OnBefore != nil {
			s.beforeHooks = append(s.beforeHooks, s.hooks[i])
		}
	}
	s.afterHookMaxIdx = len(s.afterHooks) - 1
}

// compile define a API, that Compile/Initial HostServer.
// It's only run once and should be call on s.Serve(...) before
// Gw framework s.Serve(...) will be automatic call this API.
func (s *HostServer) compile() {
	s.locker.Lock()
	defer s.locker.Unlock()
	if s.state > 0 {
		return
	}
	// All of app server initial AT here.
	initial(s)
	registerApps(s)
	prepareHooks(s)
	s.state++
}

// GetRouters returns has registered routers on the Server.
func (s *HostServer) GetRouters() []RouterInfo {
	s.compile()
	var routerInfos = make([]RouterInfo, len(s.router.routerInfos))
	copy(routerInfos, s.router.routerInfos)
	return routerInfos
}

// PrintRouterInfo ...
func (s *HostServer) PrintRouterInfo() {
	rts := s.conf.Settings.GwFramework.PrintRouterInfo
	if !rts.Disabled {
		logger.NewLine(2)
		logger.Info("%s ", rts.Title)
		logger.Info("=======================")
		var routers = s.GetRouters()
		for _, r := range routers {
			logger.Info("%s", r.String())
		}
	}
}

// Serve start the Server.
func (s *HostServer) Serve() {
	s.compile()
	// signal watch.
	sigs := make(chan os.Signal, 1)
	handler := s.options.StartBeforeHandler
	if handler != nil {
		err := s.options.StartBeforeHandler(s)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}

	s.PrintRouterInfo()

	logger.NewLine(2)
	logger.Info("Service Information")
	logger.Info("=======================")
	logger.Info(" Name: %s", s.conf.Service.Name)
	logger.Info(" Version: %s", s.conf.Service.Version)
	logger.Info(" Remarks: %s", s.conf.Service.Remarks)
	logger.NewLine(2)
	logger.Info(" Serving HTTP on: %s", s.options.Addr)

	logger.ResetLogFormatter()
	var err error
	go func() {
		err = s.router.server.Run(s.options.Addr)
	}()
	// TODO(Ocean): has no better solution that can be waiting for gin.Serve() completed with non-block state?
	time.Sleep(time.Second * 2)
	if err != nil {
		panic(fmt.Errorf("start server fail. err: %v", err))
	}
	s.serverStartDone <- struct{}{}
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
	s.serverExitSignal <- struct{}{}
	logger.Info("Shutdown server: %s, Addr: %s", s.options.Name, s.options.Addr)
}
