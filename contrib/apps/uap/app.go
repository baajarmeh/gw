package uap

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/conf"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"github.com/oceanho/gw/contrib/apps/uap/gwImpls"
	"github.com/oceanho/gw/contrib/apps/uap/restAPIs"
	"github.com/oceanho/gw/logger"
)

var (
	AksDecorator     = gw.NewPermAllDecorator("Aks")
	RoleDecorator    = gw.NewPermAllDecorator("Role")
	UserDecorator    = gw.NewPermAllDecorator("User")
	TenancyDecorator = gw.NewPermAllDecorator("Tenant")
)

type App struct {
	uap   conf.Uap
	store gw.IStore
}

func New() App {
	return App{}
}

func (App) New(store gw.IStore) App {
	return App{
		store: store,
	}
}

func (a App) Name() string {
	return "gw.uap"
}

func (a App) Router() string {
	return "uap"
}

func (a App) Register(router *gw.RouterGroup) {
	router.RegisterRestAPIs(&restAPIs.User{})
}

func (a App) Use(opt *gw.ServerOption) {
	use(opt)
}

func (a App) Migrate(state gw.ServerState) {
	dbModel.Migrate(state)
	state.DIProvider().Register(a)
}

func (a App) OnStart(state gw.ServerState) {
	initial(state)
}

func (a App) OnShutDown(state gw.ServerState) {

}

// helpers
func use(opt *gw.ServerOption) {
	opt.AuthManagerHandler = func(state gw.ServerState) gw.IAuthManager {
		return gwImpls.DefaultAuthManager(state)
	}
	opt.UserManagerHandler = func(state gw.ServerState) gw.IUserManager {
		return gwImpls.DefaultUserManager(state)
	}
	opt.PermissionManagerHandler = func(state gw.ServerState) gw.IPermissionManager {
		return gwImpls.DefaultPermissionManager(state)
	}
	opt.SessionStateManager = func(state gw.ServerState) gw.ISessionStateManager {
		return gwImpls.DefaultSessionManager(state)
	}
}

func initial(state gw.ServerState) {
	initPerms(state)
	initUsers(state)
}

// initial permission
func initPerms(state gw.ServerState) {
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
func initUsers(state gw.ServerState) {
	var uapCnf = conf.GetUAP(state.ApplicationConfig())
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
