package gwExtension

import (
	"github.com/oceanho/gw"
)

type AuthManager struct {
}

func (a AuthManager) Login(store gw.Store, passport, secret string) (*gw.User, error) {
	panic("implement me")
}

func (a AuthManager) Logout(store gw.Store, user *gw.User) bool {
	panic("implement me")
}
