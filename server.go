package gw

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"github.com/oceanho/gw/utils/secure"
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

// ServerOption represents a Server Options.
type ServerOption struct {
	Addr                     string
	Name                     string
	Restart                  string
	Prefix                   string
	PluginDir                string
	PluginSymbolName         string
	PluginSymbolSuffix       string
	StartHandlers            []ServerHandler
	ShutDownHandlers         []ServerHandler
	BackendStoreHandler      func(cnf *conf.ApplicationConfig) IStore
	AppConfigHandler         func(cnf *conf.BootConfig) *conf.ApplicationConfig
	Crypto                   func(conf *conf.ApplicationConfig) ICrypto
	IDGeneratorHandler       func(conf *conf.ApplicationConfig) IdentifierGenerator
	UserManagerHandler       func(state *ServerState) IUserManager
	DIProviderHandler        func(state *ServerState) IDIProvider
	AppManagerHandler        func(state *ServerState) IAppManager
	AuthManagerHandler       AuthManagerHandler
	AuthParamResolvers       []IAuthParamResolver
	AuthParamCheckerHandler  func(state *ServerState) IAuthParamChecker
	PermissionManagerHandler PermissionManagerHandler
	PermissionCheckerHandler func(state *ServerState) IPermissionChecker
	StoreDbSetupHandler      StoreDbSetupHandler
	SessionStateManager      SessionStateHandler
	StoreCacheSetupHandler   StoreCacheSetupHandler
	DbOpProcessor            *DbOpProcessor
	EventManagerHandler      func(state *ServerState) IEventManager
	RespBodyBuildFunc        RespBodyBuildFunc
	isTester                 bool
	cnf                      *conf.ApplicationConfig
	bcs                      *conf.BootConfig
}

type ServerHandler func(s *HostServer) error

// HostServer represents a Host Server.
type HostServer struct {
	Name                   string
	Store                  IStore
	Hash                   ICryptoHash
	Protect                ICryptoProtect
	PasswordSigner         IPasswordSigner
	AuthManager            IAuthManager
	AuthParamResolvers     []IAuthParamResolver
	AuthParamChecker       IAuthParamChecker
	SessionStateManager    ISessionStateManager
	SessionSidCreationFunc func(param AuthParameter) string
	PermissionChecker      IPermissionChecker
	PermissionManager      IPermissionManager
	UserManager            IUserManager
	AppManager             IAppManager
	IDGenerator            IdentifierGenerator
	DIProvider             IDIProvider
	EventManager           IEventManager
	DbOpProcessor          *DbOpProcessor
	RespBodyBuildFunc      RespBodyBuildFunc
	state                  int
	locker                 sync.Mutex
	options                *ServerOption
	router                 *Router
	apps                   map[string]internalApp
	conf                   *conf.ApplicationConfig
	httpErrHandlers        map[int][]ErrorHandler
	hooks                  []*Hook
	beforeHooks            []*Hook
	afterHooks             []*Hook
	afterHookMaxIdx        int
	authParamValidators    map[string]*regexp.Regexp
	storeDbSetupHandler    StoreDbSetupHandler
	storeCacheSetupHandler StoreCacheSetupHandler
	quit                   chan bool
	serverExitSignal       chan struct{}
	serverStartDone        chan struct{}
	serverHandlers         map[string][]func(state *ServerState)
}

// ServerState represents a Server state context object.
type ServerState struct {
	s *HostServer
}

func NewServerState(server *HostServer) *ServerState {
	server.compile()
	return &ServerState{
		s: server,
	}
}

func (ss *ServerState) Store() IStore {
	return ss.s.Store
}

func (ss *ServerState) CryptoHash() ICryptoHash {
	return ss.s.Hash
}

func (ss *ServerState) CryptoProtect() ICryptoProtect {
	return ss.s.Protect
}

func (ss *ServerState) PasswordSigner() IPasswordSigner {
	return ss.s.PasswordSigner
}

func (ss *ServerState) AuthManager() IAuthManager {
	return ss.s.AuthManager
}

func (ss *ServerState) AuthParamResolvers() []IAuthParamResolver {
	return ss.s.AuthParamResolvers
}

func (ss *ServerState) DbOpProcessor() *DbOpProcessor {
	return ss.s.DbOpProcessor
}

func (ss *ServerState) AuthParamChecker() IAuthParamChecker {
	return ss.s.AuthParamChecker
}

func (ss *ServerState) EventManager() IEventManager {
	return ss.s.EventManager
}

func (ss *ServerState) PermissionChecker() IPermissionChecker {
	return ss.s.PermissionChecker
}

func (ss *ServerState) SessionStateManager() ISessionStateManager {
	return ss.s.SessionStateManager
}

func (ss *ServerState) PermissionManager() IPermissionManager {
	return ss.s.PermissionManager
}

func (ss *ServerState) UserManager() IUserManager {
	return ss.s.UserManager
}

func (ss *ServerState) ServerOptions() *ServerOption {
	return ss.s.options
}

func (ss *ServerState) AppManager() IAppManager {
	return ss.s.AppManager
}

func (ss *ServerState) ApplicationConfig() *conf.ApplicationConfig {
	return ss.s.conf
}

func (ss *ServerState) IDGenerator() IdentifierGenerator {
	return ss.s.IDGenerator
}

func (ss *ServerState) DI() IDIProvider {
	return ss.s.DIProvider
}

func (ss *ServerState) RespBodyBuildFunc() RespBodyBuildFunc {
	return ss.s.RespBodyBuildFunc
}

var (
	appDefaultAddr               = ":8080"
	appDefaultName               = "gw.app"
	appDefaultRestart            = "always"
	appDefaultPrefix             = "/api/v1"
	appDefaultPluginSymbolName   = "AppPlugin"
	appDefaultPluginSymbolSuffix = ".so"
	appDefaultBackendHandler     = func(cnf *conf.ApplicationConfig) IStore {
		return DefaultBackend(cnf)
	}
	appDefaultAppConfigHandler = func(cnf *conf.BootConfig) *conf.ApplicationConfig {
		return conf.NewConfigWithBootConfig(cnf)
	}
	appDefaultStoreDbSetupHandler = func(c *Context, db *gorm.DB) *gorm.DB {
		return db.Set(gwDbContextKey, c)
	}
	appDefaultStoreCacheSetupHandler = func(c *Context, client *redis.Client, user User) *redis.Client {
		return client
	}
	internLogFormatter = "[$prefix-$level] $msg\n"
)

var (
	servers map[string]*internalHostServer
)

type internalHostServer struct {
	Server *HostServer
	State  *ServerState
}

func (hss *internalHostServer) SetState(state *ServerState) {
	hss.State = state
}

func init() {
	servers = make(map[string]*internalHostServer)
	logger.SetLogFormatter(internLogFormatter)
}

// NewServerOption returns a *ServerOption with bcs.
func NewServerOption(bcs *conf.BootConfig) *ServerOption {
	cnf := &ServerOption{
		Addr:                   appDefaultAddr,
		Name:                   appDefaultName,
		Restart:                appDefaultRestart,
		Prefix:                 appDefaultPrefix,
		AppConfigHandler:       appDefaultAppConfigHandler,
		PluginSymbolName:       appDefaultPluginSymbolName,
		PluginSymbolSuffix:     appDefaultPluginSymbolSuffix,
		StartHandlers:          make([]ServerHandler, 0, 4),
		ShutDownHandlers:       make([]ServerHandler, 0, 4),
		BackendStoreHandler:    appDefaultBackendHandler,
		StoreDbSetupHandler:    appDefaultStoreDbSetupHandler,
		StoreCacheSetupHandler: appDefaultStoreCacheSetupHandler,
		IDGeneratorHandler: func(conf *conf.ApplicationConfig) IdentifierGenerator {
			return DefaultIdentifierGenerator()
		},
		AppManagerHandler: func(state *ServerState) IAppManager {
			return EmptyAppManager(state)
		},
		UserManagerHandler: func(state *ServerState) IUserManager {
			return DefaultUserManager(state)
		},
		AuthManagerHandler: func(state *ServerState) IAuthManager {
			return DefaultAuthManager(state)
		},
		AuthParamResolvers: make([]IAuthParamResolver, 0, 3),
		AuthParamCheckerHandler: func(state *ServerState) IAuthParamChecker {
			return DefaultAuthParamChecker(state)
		},
		PermissionCheckerHandler: func(state *ServerState) IPermissionChecker {
			return DefaultPassPermissionChecker(state)
		},
		PermissionManagerHandler: func(state *ServerState) IPermissionManager {
			return DefaultPermissionManager(state)
		},
		SessionStateManager: func(state *ServerState) ISessionStateManager {
			return DefaultSessionStateManager(state)
		},
		Crypto: func(conf *conf.ApplicationConfig) ICrypto {
			c := conf.Security.Crypto
			return DefaultCrypto(c.Protect.Secret, c.Hash.Salt)
		},
		DIProviderHandler: func(state *ServerState) IDIProvider {
			return DefaultDIProvider(state)
		},
		EventManagerHandler: func(state *ServerState) IEventManager {
			return DefaultEventManager(state)
		},
		DbOpProcessor:     NewDbOpProcessor(),
		RespBodyBuildFunc: DefaultRespBodyBuildFunc,
		bcs:               bcs,
		isTester:          false,
	}
	return cnf
}

// DefaultServer returns a default HostServer(the server instance's bcs,svr are default config items.)
func DefaultServer() *HostServer {
	bcs := conf.DefaultBootConfig()
	return NewServerWithOption(NewServerOption(bcs))
}

// NewTesterServer returns a default HostServer(the server instance's bcs,svr are default config items.)
func NewTesterServer(name string) *HostServer {
	bcs := conf.DefaultBootConfig()
	opts := NewServerOption(bcs)
	opts.Name = name
	opts.isTester = true
	server := NewServerWithOption(opts)
	server.compile()
	return server
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

// New return a  Server with ServerOptions.
func NewServerWithOption(sopt *ServerOption) *HostServer {
	server, ok := servers[sopt.Name]
	if ok {
		logger.Warn("duplicated server, name: %s", sopt.Name)
		return server.Server
	}
	serverInstance := &HostServer{
		options:             sopt,
		hooks:               make([]*Hook, 0),
		beforeHooks:         make([]*Hook, 0),
		afterHooks:          make([]*Hook, 0),
		apps:                make(map[string]internalApp),
		httpErrHandlers:     make(map[int][]ErrorHandler),
		authParamValidators: make(map[string]*regexp.Regexp),
		serverExitSignal:    make(chan struct{}, 1),
		serverStartDone:     make(chan struct{}, 1),
		quit:                make(chan bool, 1),
	}
	servers[sopt.Name] = &internalHostServer{
		State:  nil,
		Server: serverInstance,
	}
	return serverInstance
}

func (s *HostServer) State() *ServerState {
	return servers[s.options.Name].State
}

func (s *HostServer) Use(f func(s *HostServer) *HostServer) *HostServer {
	return f(s)
}

// AddHook register a global http handler API into the server.
func (s *HostServer) AddHook(handlers ...*Hook) *HostServer {
	if len(handlers) == 0 {
		return s
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
	return s
}

// DeleteHook delete a global hook from has registered before hooks.
func (s *HostServer) DeleteHook(name string) *HostServer {
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
	return s
}

// ReplaceHook replace a global hook from has registered before hooks.
func (s *HostServer) ReplaceHook(name string, hookHandler *Hook) *HostServer {
	s.locker.Lock()
	defer s.locker.Unlock()
	for i := 0; i < len(s.hooks); i++ {
		if s.hooks[i].Name == name {
			s.hooks[i] = hookHandler
			break
		}
	}
	return s
}

// HandleError register a global http error handler API into the server.
func (s *HostServer) HandleErrors(handler ErrorHandler, httpStatus ...int) *HostServer {
	for _, c := range httpStatus {
		s.HandleError(c, handler)
	}
	return s
}

// HandleError register a global http error handler API into the server.
func (s *HostServer) HandleError(httpStatus int, handlers ...ErrorHandler) *HostServer {
	if len(handlers) == 0 {
		return s
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
	return s
}

// OnStart define a API that called when the server start before(Serve API before).
func (s *HostServer) OnStart(handlers ...func(state *ServerState)) *HostServer {
	s.locker.Lock()
	defer s.locker.Unlock()
	if len(s.serverHandlers) == 0 {
		s.serverHandlers = make(map[string][]func(state *ServerState))
	}
	if len(s.serverHandlers["start"]) == 0 {
		s.serverHandlers["start"] = make([]func(*ServerState), 0, 8)
	}
	s.serverHandlers["start"] = append(s.serverHandlers["start"], handlers...)
	return s
}

// OnShutDown define a API that called when the server shutdown.
func (s *HostServer) OnShutDown(handlers ...func(state *ServerState)) *HostServer {
	s.locker.Lock()
	defer s.locker.Unlock()
	if len(s.serverHandlers) == 0 {
		s.serverHandlers = make(map[string][]func(state *ServerState))
	}
	if len(s.serverHandlers["shutdown"]) == 0 {
		s.serverHandlers["shutdown"] = make([]func(*ServerState), 0, 8)
	}
	s.serverHandlers["shutdown"] = append(s.serverHandlers["shutdown"], handlers...)
	return s
}

// Register register a app instances into the server.
func (s *HostServer) Register(apps ...App) *HostServer {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, app := range apps {
		appName := app.Meta().Name
		if _, ok := s.apps[appName]; !ok {
			s.apps[appName] = internalApp{
				instance:    app,
				isPatchOnly: false,
			}
		}
	}
	return s
}

// Patch patch up for HostServer with apps.
// It's Usage scenarios are need change server Options with other gw.App but not need that's router features.
//
// Example:
//
// Your have two gw App
//
// 1. UAP(your app one), It's implementation User/Permission features
//
// 2. DevOps(your app two), It's implementation DevOps features, It's dependencies on UAP, but the UAP's router not necessary on DevOps service.
//
// Now, We can separate deployment of UAP and DevOps
//
// And we should be patch up UAP into DevOps app.
func (s *HostServer) Patch(apps ...App) *HostServer {
	if len(apps) == 0 {
		return s
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, app := range apps {
		app := app
		appName := app.Meta().Name
		if _, ok := s.apps[appName]; !ok {
			s.apps[appName] = internalApp{
				instance:    app,
				isPatchOnly: true,
			}
		}
	}
	return s
}

// RegisterByPluginDir register app instances into the server with gw plugin mode.
func (s *HostServer) RegisterByPluginDir(dirs ...string) *HostServer {
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
	return s
}

func initialConfig(s *HostServer) {
	//
	// Before Server start, Must initial all of Server Options
	// There are options may be changed by Custom app instance's .Use(...) APIs.
	cnf := s.options.AppConfigHandler(s.options.bcs)
	cnf.Compile()
	s.conf = cnf
	s.options.cnf = cnf
}

func initialServer(s *HostServer) *ServerState {
	var cnf = s.conf
	crypto := s.options.Crypto(cnf)
	s.Hash = crypto.Hash()
	s.Name = s.options.Name
	s.Protect = crypto.Protect()
	s.PasswordSigner = crypto.Password()
	s.Store = s.options.BackendStoreHandler(cnf)
	s.IDGenerator = s.options.IDGeneratorHandler(cnf)
	if s.RespBodyBuildFunc == nil {
		s.RespBodyBuildFunc = s.options.RespBodyBuildFunc
	}

	state := &ServerState{
		s: s,
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
	if len(s.options.AuthParamResolvers) == 0 {
		s.options.AuthParamResolvers = DefaultAuthParamResolver()
	}
	if s.options.AuthParamCheckerHandler == nil {
		s.options.AuthParamCheckerHandler = func(state *ServerState) IAuthParamChecker {
			return DefaultAuthParamChecker(state)
		}
	}
	// SessionSidCreationFunc
	if s.SessionSidCreationFunc == nil {
		if s.conf.Security.Auth.Session.SidGenerator == "p" {
			s.SessionSidCreationFunc = func(param AuthParameter) string {
				return param.Passport
			}
		} else if s.conf.Security.Auth.Session.SidGenerator == "p,md5" {
			s.SessionSidCreationFunc = func(param AuthParameter) string {
				return secure.Md5Str(param.Passport)
			}
		} else if s.conf.Security.Auth.Session.SidGenerator == "p,smd5" {
			s.SessionSidCreationFunc = func(param AuthParameter) string {
				return s.PasswordSigner.Sign(param.Passport)
			}
		} else {
			s.SessionSidCreationFunc = func(param AuthParameter) string {
				return s.IDGenerator.NewStrID(32)
			}
		}
	}
	if s.options.DIProviderHandler == nil {
		s.options.DIProviderHandler = func(state *ServerState) IDIProvider {
			return DefaultDIProvider(state)
		}
	}
	if s.options.EventManagerHandler == nil {
		s.options.EventManagerHandler = func(state *ServerState) IEventManager {
			return DefaultEventManager(state)
		}
	}
	if s.options.PermissionCheckerHandler == nil {
		s.options.PermissionCheckerHandler = func(state *ServerState) IPermissionChecker {
			return DefaultPassPermissionChecker(state)
		}
	}
	if s.options.AppConfigHandler == nil {
		s.options.AppManagerHandler = func(state *ServerState) IAppManager {
			return EmptyAppManager(state)
		}
	}
	s.AppManager = s.options.AppManagerHandler(state)
	s.UserManager = s.options.UserManagerHandler(state)
	s.AuthManager = s.options.AuthManagerHandler(state)
	s.AuthParamResolvers = s.options.AuthParamResolvers
	s.AuthParamChecker = s.options.AuthParamCheckerHandler(state)
	s.SessionStateManager = s.options.SessionStateManager(state)
	s.PermissionChecker = s.options.PermissionCheckerHandler(state)
	s.PermissionManager = s.options.PermissionManagerHandler(state)
	s.EventManager = s.options.EventManagerHandler(state)
	s.DbOpProcessor = s.options.DbOpProcessor
	s.DIProvider = s.options.DIProviderHandler(state)
	if !s.options.isTester {
		// gin engine.
		g := gin.New()
		// g.Use(gin.Recovery())
		g.Use(gwState(s.options.Name))

		// Auth(login/logout) API routers.
		registerBuiltinRouter(s.conf, g)

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

		// initial routes.
		httpRouter.router = httpRouter.server.Group(s.options.Prefix)
		s.router = httpRouter
	}
	if s.conf.Server.ListenAddr != "" && s.options.Addr == appDefaultAddr {
		s.options.Addr = s.conf.Server.ListenAddr
	}
	return state
}

func registerBuiltinRouter(cnf *conf.ApplicationConfig, router *gin.Engine) {
	_pprof := cnf.Service.PProf
	if _pprof.Enabled {
		pprof.Register(router, _pprof.Router)
	}
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

func useApps(s *HostServer) {
	for _, app := range s.apps {
		app.instance.Use(s.options)
	}
}

func onStarts(s *HostServer, state *ServerState) {
	for _, app := range s.apps {
		_ = s.AppManager.Create(app.instance.Meta())
		app.instance.OnStart(state)
	}
}

func registerApps(s *HostServer, state *ServerState) {
	for _, app := range s.apps {
		meta := app.instance.Meta()
		if !app.isPatchOnly {
			logger.Info("register app: %s", meta.Name)
			rg := s.router.Group(meta.Router, nil)
			app.instance.Register(rg)
		}
		// migrate
		logger.Info("migrate app: %s", meta.Name)
		app.instance.OnPrepare(state)
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
	initialConfig(s)
	useApps(s)
	state := initialServer(s)
	registerApps(s, state)
	prepareHooks(s)
	onStarts(s, state)
	servers[s.options.Name].SetState(state)
	s.state++
	go func() {
		_ = <-s.quit
		for _, handler := range s.options.ShutDownHandlers {
			err := handler(s)
			if err != nil {
				fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
			}
		}
		for _, app := range s.apps {
			app := app.instance
			app.OnShutDown(servers[s.options.Name].State)
		}
		s.serverExitSignal <- struct{}{}
		logger.Info("Shutdown server: %s, Addr: %s", s.options.Name, s.options.Addr)
	}()
}

func (s *HostServer) startHandlers() {
	if len(s.serverHandlers) == 0 {
		return
	}
	state := servers[s.Name].State
	for _, h := range s.serverHandlers["start"] {
		h(state)
	}
}

func (s *HostServer) shutDownHandlers() {
	if len(s.serverHandlers) == 0 {
		return
	}
	state := servers[s.Name].State
	for _, h := range s.serverHandlers["shutdown"] {
		h(state)
	}
}

// GetRouters returns has registered routers on the Server.
func (s *HostServer) GetRouters() []RouterInfo {
	s.compile()
	var routerInfos = make([]RouterInfo, len(s.router.routerInfos))
	copy(routerInfos, s.router.routerInfos)
	return routerInfos
}

// DisplayRouterInfo ...
func (s *HostServer) DisplayRouterInfo() {
	prtInfo := s.conf.Settings.GwFramework.PrintRouterInfo
	if !prtInfo.Disabled {
		logger.NewLine(2)
		logger.Info("%s ", prtInfo.Title)
		logger.Info("=======================")
		var routers = s.GetRouters()
		for _, r := range routers {
			logger.Info("%s", r.String())
		}
	}
}

// RegisterStartHandler start the Server.
func (s *HostServer) RegisterStartHandler(handlers ...ServerHandler) *HostServer {
	s.options.StartHandlers = append(s.options.StartHandlers, handlers...)
	return s
}

// RegisterShutDownHandler start the Server.
func (s *HostServer) RegisterShutDownHandler(handlers ...ServerHandler) *HostServer {
	s.options.ShutDownHandlers = append(s.options.ShutDownHandlers, handlers...)
	return s
}

// Serve start the Server.
func (s *HostServer) Serve() {
	s.compile()
	// signal watch.
	sigs := make(chan os.Signal, 1)
	for _, handler := range s.options.StartHandlers {
		err := handler(s)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	s.startHandlers()
	s.DisplayRouterInfo()

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
	_ = <-sigs
	s.ShutDown()
}

func (s *HostServer) ShutDown() {
	s.quit <- true
	s.shutDownHandlers()
}
