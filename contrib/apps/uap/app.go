package uap

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Api"
	"github.com/oceanho/gw/contrib/apps/uap/Config"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
	"github.com/oceanho/gw/contrib/apps/uap/Impl"
	"github.com/oceanho/gw/contrib/apps/uap/RestAPI"
	"github.com/oceanho/gw/contrib/apps/uap/Service"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
)

var (
	AksDecorator        = gw.NewPermAllDecorator("Aks")
	RoleDecorator       = gw.NewPermAllDecorator("Role")
	UserDecorator       = gw.NewPermAllDecorator("User")
	TenancyDecorator    = gw.NewPermAllDecorator("Tenant")
	CredentialDecorator = gw.NewPermAllDecorator("Credential")
)

type App struct {
	name            string
	router          string
	registerFunc    func(router *gw.RouterGroup)
	useFunc         func(option *gw.ServerOption)
	migrateFunc     func(state *gw.ServerState)
	onStartFunc     func(state *gw.ServerState)
	onShoutDownFunc func(state *gw.ServerState)
}

func New() App {
	return App{
		name:   "gw.uap",
		router: "uap",
		registerFunc: func(router *gw.RouterGroup) {
			router.RegisterRestAPIs(&RestAPI.User{})
			router.GET("credential/:id", Api.QueryCredentialById, Api.QueryCredentialByIdDecorators())
		},
		useFunc: func(option *gw.ServerOption) {
			option.AuthManagerHandler = func(state *gw.ServerState) gw.IAuthManager {
				return Impl.DefaultAuthManager(state)
			}
			option.UserManagerHandler = func(state *gw.ServerState) gw.IUserManager {
				return Impl.DefaultUserManager(state)
			}
			option.PermissionManagerHandler = func(state *gw.ServerState) gw.IPermissionManager {
				return Impl.DefaultPermissionManager(state)
			}
			option.SessionStateManager = func(state *gw.ServerState) gw.ISessionStateManager {
				return Impl.DefaultSessionManager(state)
			}
		},
		migrateFunc: func(state *gw.ServerState) {
			state.Store().GetDbStore().AutoMigrate(
				Db.User{},
				Db.UserProfile{},
				Db.Role{},
				Db.UserRoleMapping{},
				Db.Permission{},
				Db.ObjectPermission{},
				Db.Credential{},
			)
			state.DbOpProcessor().CreateBefore().Register(func(db *gorm.DB, ctx *gw.Context, model interface{}) error {
				return nil
			}, Db.Credential{})
		},
		onStartFunc: func(state *gw.ServerState) {
			// Services dependency injection
			Service.Register(state.DI())

			// TODO(OceanHo): there are may be consider initial by other tool
			//  Because of initialization with this way has some problem on a distributed Cluster.
			initPerms(state)
			initUsers(state)
		},
		onShoutDownFunc: func(state *gw.ServerState) {

		},
	}
}

func (a App) Name() string {
	return a.name
}

func (a App) Router() string {
	return a.router
}

func (a App) Register(router *gw.RouterGroup) {
	a.registerFunc(router)
}

func (a App) Use(opt *gw.ServerOption) {
	a.useFunc(opt)
}

func (a App) Migrate(state *gw.ServerState) {
	a.migrateFunc(state)
}

func (a App) OnStart(state *gw.ServerState) {
	a.onStartFunc(state)
}

func (a App) OnShutDown(state *gw.ServerState) {

}

// initial permission
func initPerms(state *gw.ServerState) {
	var perms []gw.Permission
	perms = append(perms, UserDecorator.Permissions()...)
	perms = append(perms, TenancyDecorator.Permissions()...)
	perms = append(perms, AksDecorator.Permissions()...)
	perms = append(perms, RoleDecorator.Permissions()...)
	err := state.PermissionManager().Create("uap", perms...)
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
		user.TenantId = usr.TenantId
		user.Passport = usr.Passport
		user.Secret = passwordSigner.Sign(usr.Secret)
		user.UserType = usr.UserType
		err := userManager.Create(&user)
		if err != nil && err != gw.ErrorUserHasExists {
			panic(fmt.Sprintf("uap -> initSystemAdministrator fail, err: %v", err))
		}
	}
}
