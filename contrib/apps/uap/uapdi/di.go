package uapdi

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/reposities"
	"github.com/oceanho/gw/contrib/apps/uap/services"
	"reflect"
)

var (
	serviceTyper = reflect.TypeOf(Service{})
)

func Register(di gw.IDIProvider) {
	services.Register(di)
	reposities.Register(di)
	di.RegisterWithTyper(serviceTyper)
}

type Service struct {
	UserService services.IUserService
}

func (service Service) New(userService services.IUserService) Service {
	service.UserService = userService
	return service
}

func Services(ctx *gw.Context) Service {
	return ctx.ResolveByTyper(serviceTyper).(Service)
}
