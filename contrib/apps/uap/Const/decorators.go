package Const

import (
	"github.com/oceanho/gw"
)

const (
	DecoratorCatalog = "gw/uap-check-res-op-perm-decorator"
)

var (
	AksDecorator        = gw.NewPermAllDecorator("Aks")
	RoleDecorator       = gw.NewPermAllDecorator("Role")
	RoleGlobalDecorator = gw.NewBeforeCustomFuncDecorator(func(ctx *gw.Context) (status int, err error, payload interface{}) {
		return 0, nil, nil
	})
	UserDecorator       = gw.NewPermAllDecorator("User")
	TenancyDecorator    = gw.NewPermAllDecorator("Tenant")
	CredentialDecorator = gw.NewPermAllDecorator("Credential")
)
