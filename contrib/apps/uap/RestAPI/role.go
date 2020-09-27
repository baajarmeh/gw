package RestAPI

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Dto"
	"github.com/oceanho/gw/contrib/apps/uap/Service"
)

type Role struct {
}

func (u Role) Name() string {
	return "role"
}

//
// APIs
//
// var roleDto dto.RoleDto
func (u Role) Get(ctx *gw.Context) {

}

//
//
func (u Role) Detail(ctx *gw.Context) {

}

//func (u Role) OnGetBefore() gw.Decorator {
//}

func (u Role) Post(ctx *gw.Context) {

}

// Put, Creation & decorators
func (u Role) Put(ctx *gw.Context) {
	var model Dto.Role
	if err := ctx.Bind(&model); err != nil {
		return
	}
	err := Service.Services(ctx).RoleSvc.Create(&model)
	ctx.JSON(err, nil)
}

// Delete, Deletion & decorators
func (u Role) Delete(ctx *gw.Context) {

}

// QueryList, Query Pager data & decorators
func (u Role) QueryList(ctx *gw.Context) {
	var model gw.QueryExpr
	if err := ctx.Bind(&model); err != nil {
		return
	}
	//var svc Service.Service
	//ctx.Resolve(&svc)
	err, result := Service.Services(ctx).RoleSvc.QueryList(model)
	ctx.JSON(err, result)
}
