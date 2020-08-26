package reposities

import "github.com/oceanho/gw"

func Register(di gw.IDIProvider) {
	di.Register(UserRepositoryImpl{})
}
