package uap

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/uap/api"
	"github.com/oceanho/gw/web/apps/uap/entities"
)

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (u App) Name() string {
	return "gw.uap"
}

func (u App) Router() string {
	return "uap"
}

func (u App) Register(router *gw.RouteGroup) {
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

	router.GET("perms/role/get", api.GetRolePerms)
	router.GET("perms/role/query", api.QueryRolePerms)
	router.POST("perms/role/create", api.CreateRolePerms)
	router.POST("perms/role/modify", api.ModifyRolePerms)
	router.POST("perms/role/destroy", api.DeleteRolePerms)

	router.GET("perms/user/get", api.GetUserPerms)
	router.GET("perms/user/query", api.QueryUserPerms)
	router.POST("perms/user/create", api.CreateUserPerms)
	router.POST("perms/user/modify", api.ModifyUserPerms)
	router.POST("perms/user/destroy", api.DeleteUserPerms)
}

func (u App) Migrate(store gw.Store) {
	db := store.GetDbStore()
	db.AutoMigrate(&entities.User{}, &entities.Profile{})
}
