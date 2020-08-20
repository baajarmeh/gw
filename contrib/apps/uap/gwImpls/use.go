package gwImpls

import "github.com/oceanho/gw"

func Use(opt *gw.ServerOption) {
	opt.AuthManagerHandler = func(state gw.ServerState) gw.IAuthManager {
		return DefaultAuthManager(state)
	}
	opt.UserManagerHandler = func(state gw.ServerState) gw.IUserManager {
		return DefaultUserManager(state)
	}
	opt.PermissionManagerHandler = func(state gw.ServerState) gw.IPermissionManager {
		return DefaultPermissionManager(state)
	}
}
