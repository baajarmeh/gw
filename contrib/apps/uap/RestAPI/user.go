package RestAPI

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Service"
)

type User struct {
}

func (u User) Name() string {
	return "user"
}

//
// APIs
//
// var userDto dto.UserDto
func (u User) Get(ctx *gw.Context) {
	user := ctx.User()
	svc := Service.Services(ctx)
	svc.UserManager.Query(user.TenantID)
}

//
//
func (u User) Detail(ctx *gw.Context) {
	userManager := Service.Services(ctx).UserManager
	_, _ = userManager.QueryByUser("admin", "admin@123456")
}

//func (u User) OnGetBefore() gw.Decorator {
//}

func (u User) Post(ctx *gw.Context) {

}

// Put, Creation & decorators
func (u User) Put(ctx *gw.Context) {
}

// Delete, Deletion & decorators
func (u User) Delete(ctx *gw.Context) {

}

// QueryList, Query Pager data & decorators
func (u User) QueryList(ctx *gw.Context) {

}
