package Service

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Impl"
	"reflect"
)

type Service struct {
	gw.BuiltinComponent
	UserManager       gw.IUserManager
	RoleSvc           IRoleService
	CredentialService ICredentialService
}

// DI
func (s Service) New(userManager gw.IUserManager, builtin gw.BuiltinComponent,
	service ICredentialService, roleSvc IRoleService) Service {
	s.BuiltinComponent = builtin
	s.CredentialService = service
	s.UserManager = userManager
	s.RoleSvc = roleSvc
	return s
}

var serviceTyper = reflect.TypeOf(Service{})

func Services(ctx *gw.Context) Service {
	return ctx.ResolveByTyper(serviceTyper).(Service)
}

func Register(di gw.IDIProvider) {
	registerServices(di)
	di.RegisterWithTyper(serviceTyper)
}

func registerServices(di gw.IDIProvider) {
	di.Register(Impl.AppManager{}, Impl.UserManager{}, RoleService{})
	di.Register(DefaultCredentialProtectServiceImpl{}, DefaultCredentialServiceImpl{})
}
