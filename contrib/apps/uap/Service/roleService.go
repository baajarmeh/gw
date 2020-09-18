package Service

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
	"github.com/oceanho/gw/contrib/apps/uap/Dto"
)

type IRoleService interface {
	Create(dto *Dto.Role) error
	Modify(dto *Dto.Role) error
	Delete(id uint64) error
	Query(id uint64) (error, Dto.Role)
	QueryList(expr gw.QueryExpr) (error, RolePagerResult)
	QueryListByTenant(tenantId uint64, expr gw.QueryExpr) (error, RolePagerResult)
	QueryListByUser(userId uint64, expr gw.QueryExpr) (error, RolePagerResult)
}

type RolePagerResult struct {
	gw.PagerResult
	Data []Dto.Role `json:"data"`
}

type RoleService struct {
	gw.BuiltinComponent
}

// DI

func (rs RoleService) New(bc gw.BuiltinComponent) IRoleService {
	rs.BuiltinComponent = bc
	return rs
}

func (rs RoleService) Create(dto *Dto.Role) error {
	var role Db.Role
	role.Name = dto.Name
	role.Descriptor = dto.Descriptor
	if rs.User.IsTenancy() {
		role.TenantID = rs.User.ID
	} else if rs.User.IsUser() {
		role.TenantID = rs.User.TenantID
	}
	store := rs.Store.GetDbStore()
	return store.Create(&role).Error
}

func (rs RoleService) Modify(dto *Dto.Role) error {
	panic("implement me")
}

func (rs RoleService) Delete(id uint64) error {
	panic("implement me")
}

func (rs RoleService) Query(id uint64) (error, Dto.Role) {
	panic("implement me")
}

func (rs RoleService) QueryList(expr gw.QueryExpr) (error, RolePagerResult) {
	var result RolePagerResult
	result.PagerExpr = expr.PagerExpr
	result.Data = make([]Dto.Role, 0, expr.PageSize)
	roleDb := rs.Store.GetDbStore().Model(&Db.Role{})
	err := roleDb.Count(&result.Total).Offset(expr.PageOffset()).Limit(expr.PageSize).Find(&result.Data).Error
	return err, result
}

func (rs RoleService) QueryListByTenant(tenantId uint64, expr gw.QueryExpr) (error, RolePagerResult) {
	panic("implement me")
}

func (rs RoleService) QueryListByUser(userId uint64, expr gw.QueryExpr) (error, RolePagerResult) {
	panic("implement me")
}
