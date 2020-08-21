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

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (a App) Name() string {
	return "gw.uap"
}

func (a App) Router() string {
	return "uap"
}

func (a App) Register(router *gw.RouterGroup) {
	restAPIs.Register(router)
}

func (a App) Use(opt *gw.ServerOption) {
	use(opt)
}

func (a App) Migrate(state gw.ServerState) {
	dbModel.Migrate(state)
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
}

func initial(state gw.ServerState) {
	initSystemAdministrator(state)
}

// initial system administrator
func initSystemAdministrator(state gw.ServerState) {
	var cnfUsers []conf.User
	var cnf = state.ApplicationConfig()
	err := cnf.ParseCustomPathTo("gwcontrib.uap.initialUsers", &cnfUsers)
	_ = cnf.ParseCustomPathTo("gwcontrib.uap.initialUsers", &cnfUsers)
	if err != nil {
		logger.Error("read gwcontrib.uap.initialUsers fail, err: %v", err)
		return
	}
	var userManager = state.UserManager()
	var passwordSigner = state.PasswordSigner()
	for _, u := range cnfUsers {
		usr := u
		var user gw.AuthUser
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
