package gwImpls

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"gorm.io/gorm"
	"sync"
)

type PermissionManagerImpl struct {
	_state int
	store  gw.IStore
	conf   conf.ApplicationConfig
	locker sync.Mutex
	state  gw.ServerState
}

func DefaultPermissionManager(state gw.ServerState) gw.IPermissionManager {
	return &PermissionManagerImpl{
		state: state,
		conf:  *state.ApplicationConfig(),
		store: state.Store(),
	}
}

func (p *PermissionManagerImpl) Store() *gorm.DB {
	return p.store.GetDbStoreByName(p.conf.Security.Auth.Permission.DefaultStore.Name)
}

func (p *PermissionManagerImpl) Initial() {
	p.locker.Lock()
	defer p.locker.Unlock()
	if p._state > 0 {
		return
	}
	store := p.Store()
	var tables []interface{}
	tables = append(tables, &dbModel.Permission{})
	tables = append(tables, &dbModel.PermissionMapping{})
	err := store.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}

func (p *PermissionManagerImpl) Has(user gw.User, perms ...gw.Permission) bool {
	if user.IsAuth() {
		if user.IsAdmin() {
			return true
		}
		for _, g := range user.Permissions {
			for _, pm := range perms {
				if g.Key == pm.Key {
					return true
				}
			}
		}
	}
	return false
}

func (p *PermissionManagerImpl) Create(category string, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return gw.ErrorEmptyInput
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	db := p.Store()
	tx := db.Begin()
	var registered = make(map[string]bool)
	for i := 0; i < len(perms); i++ {
		p := perms[i]
		uk := fmt.Sprintf("%d-%s-%s", p.TenantId, p.Category, p.Key)
		if registered[uk] {
			continue
		}
		var count int64
		if p.Category == "" {
			p.Category = category
		}
		err := db.Model(dbModel.Permission{}).Where("tenant_id = ? and category = ? and `key` = ?",
			p.TenantId, p.Category, p.Key).Count(&count).Error
		if err != nil {
			panic(fmt.Sprintf("check perms has exist fail, err: %v", err))
		}
		if count == 0 {
			tx.Create(p)
		}
	}
	return tx.Commit().Error
}

func (p *PermissionManagerImpl) Modify(perms ...gw.Permission) error {
	if len(perms) < 1 {
		return gw.ErrorEmptyInput
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	tx := p.Store().Begin()
	for _, p := range perms {
		tx.Updates(p)
	}
	return tx.Commit().Error
}

func (p *PermissionManagerImpl) Drop(perms ...gw.Permission) error {
	if len(perms) < 1 {
		return gw.ErrorEmptyInput
	}
	tx := p.Store().Begin()
	for _, p := range perms {
		tx.Delete(p)
	}
	return tx.Commit().Error
}

func (p *PermissionManagerImpl) Query(tenantId uint64, category string, expr gw.PagerExpr) (
	total int64, result []gw.Permission, error error) {
	error = p.Store().Where("tenant_id = ? and category = ?",
		tenantId, category).Count(&total).Offset(expr.PageOffset()).Limit(expr.PageSize).Scan(result).Error
	return
}

func (p *PermissionManagerImpl) QueryByUser(tenantId, userId uint64, expr gw.PagerExpr) (
	total int64, result []gw.Permission, error error) {
	var perm dbModel.Permission
	var objPerm dbModel.PermissionMapping
	var sql = fmt.Sprintf(" from %s t1 inner join %s t2 on t1.id = t2.permission_id "+
		" where t2.tenant_id=%d and t2.object_id=%d and t2.type=%d",
		perm.TableName(), objPerm.TableName(), tenantId, userId, dbModel.UserPermission)
	var countSql = fmt.Sprintf("select count(t1.Id) as total %s", sql)
	var dataSql = fmt.Sprintf("select t1.* %s limit %d offset %d", sql, expr.PageSize, expr.PageOffset())
	tx := p.Store().Begin()
	err := tx.Raw(countSql).Scan(&total).Error
	if err != nil {
		tx.Rollback()
		return total, nil, err
	}
	err = tx.Raw(dataSql).Scan(&result).Error
	if err != nil {
		tx.Rollback()
		return 0, nil, err
	}
	err = tx.Commit().Error
	return total, result, err
}

func (p *PermissionManagerImpl) GrantToUser(uid uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := p.Store().Begin()
	for _, p := range perms {
		var pm dbModel.PermissionMapping
		pm.PermissionID = p.ID
		pm.TenantId = p.TenantId
		pm.Type = dbModel.UserPermission
		pm.ObjectID = uid
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (p *PermissionManagerImpl) GrantToRole(roleId uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := p.Store().Begin()
	for _, p := range perms {
		var pm dbModel.PermissionMapping
		pm.PermissionID = p.ID
		pm.TenantId = p.TenantId
		pm.Type = dbModel.RolePermission
		pm.ObjectID = roleId
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (p *PermissionManagerImpl) RevokeFromUser(uid uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := p.Store().Begin()
	for _, p := range perms {
		tx.Delete(dbModel.PermissionMapping{},
			"object_id = ? and tenant_id = ? and permission_id = ? and type = ?", uid, p.TenantId, p.ID, dbModel.UserPermission)
	}
	return tx.Commit().Error
}

func (p *PermissionManagerImpl) RevokeFromRole(roleId uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := p.Store().Begin()
	for _, p := range perms {
		tx.Delete(dbModel.PermissionMapping{},
			"object_id = ? and tenant_id = ? and permission_id = ? and type = ?", roleId, p.TenantId, p.ID, dbModel.RolePermission)
	}
	return tx.Commit().Error
}
