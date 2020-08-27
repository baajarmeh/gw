package restAPIs

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/constants"
	"github.com/oceanho/gw/contrib/apps/uap/dto"
	"github.com/oceanho/gw/contrib/apps/uap/uapdi"
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
	var id uint64
	if ctx.MustGetParamIDUint64(&id) != nil {
		// If binding fail, error message has sent by GW framework.
		// Here returns then done.
		return
	}
	var services = uapdi.Services(ctx)
	dto, err := services.UserService.GetById(id)
	ctx.JSON(err, dto)
}

//func (u User) OnGetBefore() gw.Decorator {
//}

func (u User) Post(ctx *gw.Context) {

}

func (u User) OnPostBefore() gw.Decorator {
	return constants.UserDecorators.Modification()
}

// Put, Creation & decorators
func (u User) Put(ctx *gw.Context) {
	var dto dto.UserDto
	if ctx.Bind(&dto) != nil {
		// If binding fail, error message has sent by GW framework.
		// Here returns then done.
		return
	}
	var services = uapdi.Services(ctx)
	err := services.UserService.Create(dto)
	ctx.JSON(err, nil)
}

func (u User) OnPutBefore() gw.Decorator {
	return constants.UserDecorators.Creation()
}

// Delete, Deletion & decorators
func (u User) Delete(ctx *gw.Context) {

}

// QueryList, Query Pager data & decorators
func (u User) QueryList(ctx *gw.Context) {

}
