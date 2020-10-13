package uap

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Api"
	"github.com/oceanho/gw/contrib/apps/uap/Config"
	"github.com/oceanho/gw/contrib/apps/uap/Const"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
	"github.com/oceanho/gw/contrib/apps/uap/GwImpl"
	"github.com/oceanho/gw/contrib/apps/uap/RestAPI"
	"github.com/oceanho/gw/contrib/apps/uap/Service"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"reflect"
)

const (
	appKey    = "gw.uap"
	appName   = "gw.uap"
	appRouter = "oceanho/gw-uap"
)

type App struct {
	key            string
	name           string
	router         string
	registerFunc   func(router *gw.RouterGroup)
	useFunc        func(option *gw.ServerOption)
	onPrepareFunc  func(state *gw.ServerState)
	onStartFunc    func(state *gw.ServerState)
	onShutDownFunc func(state *gw.ServerState)
}

var dbUserTableTyper = reflect.TypeOf(Db.User{})

func New() App {
	return App{
		key:    appKey,
		name:   appName,
		router: appRouter,
		registerFunc: func(router *gw.RouterGroup) {
			router.RegisterRestAPIs(&RestAPI.User{}, &RestAPI.Role{})
			router.GET("credential/:id", Api.QueryCredentialById, Api.QueryCredentialByIdDecorators())

			router.GET("menu/i/:id", Api.GetMenu)
			router.GET("menu/n/:name", Api.GetMenuByName)
			router.GET("menu/pageList", Api.QueryMenuPageList)
			router.GET("menu/search", Api.SearchMenuPageList)
			router.POST("menu/batchCreate", Api.BatchCreateMenu)
			router.POST("menu/create", Api.CreateMenu)
			router.POST("menu/modify", Api.ModifyMenu)
		},
		useFunc: func(option *gw.ServerOption) {
			option.AppManagerHandler = func(state *gw.ServerState) gw.IAppManager {
				return GwImpl.DefaultAppManager(state)
			}
			option.AuthManagerHandler = func(state *gw.ServerState) gw.IAuthManager {
				return GwImpl.DefaultAuthManager(state)
			}
			option.UserManagerHandler = func(state *gw.ServerState) gw.IUserManager {
				return GwImpl.DefaultUserManager(state)
			}
			option.PermissionManagerHandler = func(state *gw.ServerState) gw.IPermissionManager {
				return GwImpl.DefaultPermissionManager(state)
			}
			option.SessionStateManager = func(state *gw.ServerState) gw.ISessionStateManager {
				return GwImpl.DefaultSessionManager(state)
			}
		},
		onPrepareFunc: func(state *gw.ServerState) {
			var uapConf = Config.GetUAP(state.ApplicationConfig())
			err := state.Store().GetDbStoreByName(uapConf.Backend.Name).AutoMigrate(
				Db.App{},
				Db.Menu{},
				Db.User{},
				Db.UserProfile{},
				Db.Role{},
				Db.UserRole{},
				Db.Credential{},
				Db.Permission{},
				Db.PermissionRelation{},
				Db.UserAccessKeySecret{},
			)
			if err != nil {
				panic("migrate uap fail")
			}
			var globalFilterFunc = func(db *gorm.DB, authUser gw.User) error {
				if authUser.IsEmpty() {
					return nil
				}
				_, hasUserId := db.Statement.Schema.FieldsByName["UserID"]
				_, hasTenantId := db.Statement.Schema.FieldsByName["TenantID"]
				if authUser.IsTenancy() {
					if db.Statement.Schema.ModelType == dbUserTableTyper {
						db.Where("(id = ? or tenant_id = ?)", authUser.ID, authUser.ID)
					} else {
						if hasTenantId {
							db.Where("tenant_id = ?", authUser.ID)
						}
					}
				} else if authUser.IsUser() {
					if db.Statement.Schema.ModelType == dbUserTableTyper {
						db.Where("id = ?", authUser.ID)
					} else {
						if hasUserId {
							db.Where("user_id = ?", authUser.ID)
						}
						if hasTenantId {
							db.Where("tenant_id = ?", authUser.TenantID)
						}
					}
				}
				return nil
			}
			state.DbOpProcessor().UpdateBefore().Register(func(db *gorm.DB, ctx *gw.Context) error {
				return globalFilterFunc(db, ctx.User())
			})
			state.DbOpProcessor().DeleteBefore().Register(func(db *gorm.DB, ctx *gw.Context) error {
				return globalFilterFunc(db, ctx.User())
			})
			state.DbOpProcessor().QueryBefore().Register(func(db *gorm.DB, ctx *gw.Context) error {
				return globalFilterFunc(db, ctx.User())
			})
		},
		onStartFunc: func(state *gw.ServerState) {
			// Services dependency injection
			Service.Register(state.DI())
			// TODO(OceanHo): there are may be consider initial by other tool
			//  Because of initialization with this way has some problem on a distributed Cluster.
			initPerms(state)
			initUsers(state)
		},
		onShutDownFunc: func(state *gw.ServerState) {
		},
	}
}

func (a App) Meta() gw.AppInfo {
	return gw.AppInfo{
		Key:        a.key,
		Name:       a.name,
		Router:     a.router,
		Descriptor: "gw user account platform",
	}
}

func (a App) Register(router *gw.RouterGroup) {
	a.registerFunc(router)
}

func (a App) Use(opt *gw.ServerOption) {
	a.useFunc(opt)
}

func (a App) OnPrepare(state *gw.ServerState) {
	a.onPrepareFunc(state)
}

func (a App) OnStart(state *gw.ServerState) {
	a.onStartFunc(state)
}

func (a App) OnShutDown(state *gw.ServerState) {
	a.onShutDownFunc(state)
}

// initial permission
func initPerms(state *gw.ServerState) {
	var perms []*gw.Permission
	var appInfo = state.AppManager().QueryByName(appName)
	perms = append(perms, Const.UserPermDecorator.Permissions()...)
	perms = append(perms, Const.TenancyPermDecorator.Permissions()...)
	perms = append(perms, Const.AksPermDecorator.Permissions()...)
	perms = append(perms, Const.RolePermDecorator.Permissions()...)
	perms = append(perms, Const.CredentialPermDecorator.Permissions()...)
	gw.VisitPerms(perms, appInfo)
	err := state.PermissionManager().Create(perms...)
	if err != nil {
		logger.Error("initial permissions fail, err: %v", err)
		return
	}
}

// initial users
func initUsers(state *gw.ServerState) {
	var uapCnf = Config.GetUAP(state.ApplicationConfig())
	var userManager = state.UserManager()
	var passwordSigner = state.PasswordSigner()
	for _, u := range uapCnf.User.Users {
		usr := u
		var user gw.User
		user.Passport = usr.User
		user.TenantID = usr.TenantID
		user.Secret = passwordSigner.Sign(usr.Secret)
		user.UserType = usr.UserType
		err := userManager.Create(&user)
		if err != nil && err != gw.ErrorUserHasExists {
			logger.Error(fmt.Sprintf("uap -> initSystemAdministrator fail, err: %v", err))
		}
	}
}
