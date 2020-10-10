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
	UserSvc           IUserService
	CredentialService ICredentialService
}

// DI
func (s Service) New(userManager gw.IUserManager,
	builtin gw.BuiltinComponent,
	service ICredentialService,
	roleSvc IRoleService,
	userSvc IUserService) Service {
	s.BuiltinComponent = builtin
	s.CredentialService = service
	s.UserManager = userManager
	s.RoleSvc = roleSvc
	s.UserSvc = userSvc
	return s
}

var serviceTyper = reflect.TypeOf(Service{})

func Services(ctx *gw.Context) (s Service) {
	if e := ctx.ResolveByObjectTyper(&s); e != nil {
		panic(e)
	}
	return
}

func Register(di gw.IDIProvider) {
	registerServices(di)
	di.RegisterWithTyper(serviceTyper)
}

func registerServices(di gw.IDIProvider) {
	di.Register(
		Impl.AppManager{},
		Impl.UserManager{},
		Impl.SessionManager{},
		Impl.AuthManager{},
		Impl.PermissionManager{},
		RoleService{},
		UserService{},
		DefaultCredentialProtectServiceImpl{},
		DefaultCredentialServiceImpl{})
}
