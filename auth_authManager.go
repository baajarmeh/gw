package gw

import (
	"fmt"
	"github.com/oceanho/gw/conf"
)

type DefaultAuthManagerImpl struct {
	store IStore
	cnf   *conf.ApplicationConfig
	users map[string]*defaultUser
}

func DefaultAuthManager(state *ServerState) *DefaultAuthManagerImpl {
	var authManager DefaultAuthManagerImpl
	authManager.cnf = state.ApplicationConfig()
	authManager.store = state.Store()
	return &authManager
}

func (d *DefaultAuthManagerImpl) Login(param AuthParameter) (User, error) {
	user, ok := d.users[param.Passport]
	if ok && user.secret == param.Password {
		return user.User, nil
	}
	return EmptyUser, fmt.Errorf("user:%s not found or serect not match", param.Passport)
}

func (d *DefaultAuthManagerImpl) Logout(user User) bool {
	// Nothing to do.
	return true
}
