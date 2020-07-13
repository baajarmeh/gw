package uap

import (
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/web/apps/uap/api"
)

func init() {
}

type App struct {
}

func New() *App {
	return &App{}
}

func (u App) Name() string {
	return "oceanho.uap"
}

func (u App) BaseRouter() string {
	return "uap"
}

func (u App) Register(router *app.ApiRouteGroup) {
	router.GET("tenant/get", api.GetTenant)
	router.GET("tenant/query", api.QueryTenant)
	router.POST("tenant/create", api.CreateTenant)
	router.POST("tenant/modify", api.ModifyTenant)
	router.POST("tenant/destroy", api.DeleteTenant)

	router.GET("user/get", api.GetUser)
	router.GET("user/query", api.QueryUser)
	router.POST("user/create", api.CreateUser)
	router.POST("user/modify", api.ModifyUser)
	router.POST("user/destroy", api.DeleteUser)

	router.GET("ak/get", api.GetAK)
	router.GET("ak/query", api.QueryAK)
	router.POST("ak/create", api.CreateAK)
	router.POST("ak/modify", api.ModifyAK)
	router.POST("ak/destroy", api.DeleteAK)

	router.GET("perms/get", api.GetPermission)
	router.GET("perms/query", api.QueryPermission)
	router.POST("perms/create", api.CreatePermission)
	router.POST("perms/modify", api.ModifyPermission)
	router.POST("perms/destroy", api.DeletePermission)

	router.GET("role/get", api.GetRole)
	router.GET("role/query", api.QueryRole)
	router.POST("role/create", api.CreateRole)
	router.POST("role/modify", api.ModifyRole)
	router.POST("role/destroy", api.DeleteRole)

	router.GET("dept/get", api.GetDept)
	router.GET("dept/query", api.QueryDept)
	router.POST("dept/create", api.CreateDept)
	router.POST("dept/modify", api.ModifyDept)
	router.POST("dept/destroy", api.DeleteDept)

	router.GET("role-perms/get", api.GetRolePerms)
	router.GET("role-perms/query", api.QueryRolePerms)
	router.POST("role-perms/create", api.CreateRolePerms)
	router.POST("role-perms/modify", api.ModifyRolePerms)
	router.POST("role-perms/destroy", api.DeleteRolePerms)

	router.GET("user-perms/get", api.GetUserPerms)
	router.GET("user-perms/query", api.QueryUserPerms)
	router.POST("user-perms/create", api.CreateUserPerms)
	router.POST("user-perms/modify", api.ModifyUserPerms)
	router.POST("user-perms/destroy", api.DeleteUserPerms)
}
