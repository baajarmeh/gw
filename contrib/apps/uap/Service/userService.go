package Service

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
	"github.com/oceanho/gw/contrib/apps/uap/Dto"
)

type IUserService interface {
	Create(dto *Dto.User) error
	Modify(dto *Dto.User) error
	Delete(id uint64) error
	Query(id uint64) (error, Dto.User)
	QueryList(expr gw.QueryExpr) (error, UserPagerResult)
	QueryListByTenant(tenantId uint64, expr gw.QueryExpr) (error, UserPagerResult)
	QueryListByUser(userId uint64, expr gw.QueryExpr) (error, UserPagerResult)
}

type UserPagerResult struct {
	gw.PagerResult
	Data []Dto.User `json:"data"`
}

type UserService struct {
	gw.BuiltinComponent
}

// DI

func (rs UserService) New(bc gw.BuiltinComponent) IUserService {
	rs.BuiltinComponent = bc
	return rs
}

func (rs UserService) Create(dto *Dto.User) error {
	var user Db.User
	user.Passport = dto.Passport
	user.Secret = rs.PasswordSigner.Sign(dto.Secret)
	if rs.User.IsTenancy() {
		user.TenantID = rs.User.ID
	} else if rs.User.IsUser() {
		user.TenantID = rs.User.TenantID
	}
	store := rs.Store.GetDbStore()
	return store.Create(&user).Error
}

func (rs UserService) Modify(dto *Dto.User) error {
	panic("implement me")
}

func (rs UserService) Delete(id uint64) error {
	panic("implement me")
}

func (rs UserService) Query(id uint64) (error, Dto.User) {
	panic("implement me")
}

func (rs UserService) QueryList(expr gw.QueryExpr) (error, UserPagerResult) {
	var result UserPagerResult
	result.PagerExpr = expr.PagerExpr
	result.Data = make([]Dto.User, 0, expr.PageSize)
	roleDb := rs.Store.GetDbStore().Model(&Db.User{})
	err := roleDb.Count(&result.Total).Offset(expr.PageOffset()).Limit(expr.PageSize).Find(&result.Data).Error
	return err, result
}

func (rs UserService) QueryListByTenant(tenantId uint64, expr gw.QueryExpr) (error, UserPagerResult) {
	panic("implement me")
}

func (rs UserService) QueryListByUser(userId uint64, expr gw.QueryExpr) (error, UserPagerResult) {
	panic("implement me")
}
