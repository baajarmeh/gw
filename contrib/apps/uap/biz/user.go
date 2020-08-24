package biz

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"github.com/oceanho/gw/contrib/apps/uap/dto"
	"gorm.io/gorm"
)

type IUserRepo interface {
	CreateUser(dto dto.UserDto)
}

type UserRepo struct {
	User  *gorm.DB
	ctx   gw.Context
	store gw.IStore
}

func (UserRepo) New(store gw.IStore, ctx gw.Context) IUserRepo {
	return UserRepo{
		ctx:   ctx,
		store: store,
		User:  store.GetDbStore().Model(dbModel.User{}),
	}
}

func (u UserRepo) CreateUser(dto dto.UserDto) {
	//u.User
}
