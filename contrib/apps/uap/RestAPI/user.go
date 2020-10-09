package RestAPI

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Const"
	"github.com/oceanho/gw/contrib/apps/uap/Dto"
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

}

//
//
func (u User) Detail(ctx *gw.Context) {
	var uid uint64
	if err := ctx.MustGetUint64IDFromParam(&uid); err != nil {
		return
	}
	userManager := Service.Services(ctx).UserManager
	user, err := userManager.Query(uid)
	ctx.JSON(err, user)
}

//func (u User) OnGetBefore() gw.Decorator {
//}

func (u User) Post(ctx *gw.Context) {
	var user gw.User
	var model Dto.User
	if err := ctx.Bind(&model); err != nil {
		return
	}
	var auth = ctx.User()
	user.ID = model.ID
	user.UserType = model.UserType
	if model.UserType.IsAdmin() && !auth.UserType.IsAdmin() {
		ctx.JSON403Msg(403, Const.ErrorNonUserCannotModifyResource)
		return
	}
	if auth.IsUser() {
		user.ID = auth.ID
		user.TenantID = auth.TenantID
		user.UserType = gw.NormalUser
	} else if auth.IsTenancy() {
		user.TenantID = auth.TenantID
	}

	svc := Service.Services(ctx)
	user.Passport = model.Passport
	//user.Secret = svc.PasswordSigner.Sign(model.Secret)
	user.UserType = model.UserType

	err := svc.UserManager.Modify(user)
	ctx.JSON(err, nil)
}

// Put, Creation & decorators
func (u User) Put(ctx *gw.Context) {
	var model Dto.User
	if err := ctx.Bind(&model); err != nil {
		return
	}
	var auth = ctx.User()
	// Non-user can not be create User resource
	if auth.IsUser() {
		ctx.JSON403Msg(403, Const.ErrorNonUserCannotCreationResource)
		return
	}
	// Tenancy can not be create Admin/Tenancy type resource
	if auth.IsTenancy() && (model.UserType.IsAdmin() || model.UserType.IsTenancy()) {
		ctx.JSON403Msg(403, Const.ErrorTenancyCannotCreationAdminResource)
		return
	}
	if model.UserType == 0 {
		model.UserType = gw.NormalUser
	}
	var user gw.User
	svc := Service.Services(ctx)
	user.TenantID = auth.ID
	user.Passport = model.Passport
	user.Secret = svc.PasswordSigner.Sign(model.Secret)
	user.UserType = model.UserType

	err := svc.UserManager.Create(&user)
	ctx.JSON(err, nil)
}

// Delete, Deletion & decorators
func (u User) Delete(ctx *gw.Context) {

}

// QueryList, Query Pager data & decorators
func (u User) QueryList(ctx *gw.Context) {

}
