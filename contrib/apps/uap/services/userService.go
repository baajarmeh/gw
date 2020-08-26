package services

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"github.com/oceanho/gw/contrib/apps/uap/dto"
	"github.com/oceanho/gw/contrib/apps/uap/reposities"
	"gorm.io/gorm"
)

type IUserService interface {
	Create(dto dto.UserDto) error
	GetById(id uint64) dto.UserDto
}

type UserService struct {
	store          gw.IStore
	primaryDb      *gorm.DB
	passwordSigner gw.IPasswordSigner
	UserRepo       reposities.UserRepository
}

func (u UserService) New(store gw.IStore, passwordSigner gw.IPasswordSigner, userRepo reposities.UserRepository) IUserService {
	u.store = store
	u.UserRepo = userRepo
	u.passwordSigner = passwordSigner
	u.primaryDb = store.GetDbStore()
	return u
}

func (u UserService) Create(dto dto.UserDto) error {
	var model dbModel.User
	model.TenantId = dto.TenantId
	model.IsUser = dto.UserType.IsUser()
	model.IsTenancy = dto.UserType.IsTenancy()
	model.IsAdmin = dto.UserType.IsAdmin()
	model.Secret = u.passwordSigner.Sign(dto.Secret)
	model.IsActive = dto.IsActive
	return u.UserRepo.Create(&model)
}

func (u UserService) GetById(id uint64) dto.UserDto {
	var model dbModel.User
	var dto dto.UserDto
	model.ID = id
	u.primaryDb.First(model).Scan(&dto)
	return dto
}
