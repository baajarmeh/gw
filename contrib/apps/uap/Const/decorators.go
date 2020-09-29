package Const

import (
	"github.com/oceanho/gw"
)

const (
	DecoratorCatalog = "gw/uap-check-res-op-perm-decorator"
)

var (
	AksPermDecorator    = gw.NewPermAllDecorator("Aks")
	RolePermDecorator   = gw.NewPermAllDecorator("Role")
	RoleGlobalDecorator = gw.NewBeforeCustomFuncDecorator(func(ctx *gw.Context) (status int, err error, payload interface{}) {
		return 0, nil, nil
	})
	UserPermDecorator       = gw.NewPermAllDecorator("User")
	TenancyPermDecorator    = gw.NewPermAllDecorator("Tenant")
	CredentialPermDecorator = gw.NewPermAllDecorator("Credential")
)
