package Service

import (
	"github.com/oceanho/gw"
	"reflect"
)

type Service struct {
	CredentialService ICredentialService
}

func (s Service) New(service ICredentialService) Service {
	s.CredentialService = service
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
	di.Register(DefaultCredentialProtectServiceImpl{}, DefaultCredentialServiceImpl{})
}
