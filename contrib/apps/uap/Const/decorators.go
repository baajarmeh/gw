package Const

import (
	"github.com/oceanho/gw"
)

const (
	DecoratorCatalog = "gw/uap-check-res-op-perm-decorator"
)

func globalPermValidator(ctx *gw.Context) (status int, err error, payload interface{}) {
	// Role level resources, only allow tenancy, administrator, platform user can be operate.
	user := ctx.User()
	if !user.IsPlatformUserOrTenancy() {
		return 403, ErrorNonUserCannotCreationResource, nil
	}
	return 0, nil, nil
}

var (
	AksPermDecorator        = gw.NewPermAllDecorator("Aks")
	RolePermDecorator       = gw.NewPermAllDecorator("Role")
	RoleGlobalDecorator     = gw.NewBeforeCustomFuncDecorator(globalPermValidator)
	UserPermDecorator       = gw.NewPermAllDecorator("User")
	TenancyPermDecorator    = gw.NewPermAllDecorator("Tenant")
	CredentialPermDecorator = gw.NewPermAllDecorator("Credential")
)
